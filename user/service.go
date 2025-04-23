package user

import (
	"avarts/constants"
	"avarts/utils"
	"errors"
	"mime/multipart"
	"net/http"

	"gorm.io/gorm"
)

type Service interface {
	GetByUsername(username string) (*User, error)
	GetByID(userID uint) (*User, error)
	UpdateProfile(userID uint, updated UpdateProfileRequest) error
	UploadAvatarToS3(file *multipart.File, fileHeader *multipart.FileHeader) (*string, int, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository}
}

func (s *service) GetByUsername(username string) (*User, error) {
	return s.repository.GetByUsername(username)
}

func (s *service) GetByID(userID uint) (*User, error) {
	return s.repository.Get(userID)
}

func (s *service) UpdateProfile(userID uint, updated UpdateProfileRequest) error {
	user, err := s.repository.Get(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(constants.UserNotFound)
		}
		return err
	}

	if updated.Username != nil {
		newUsername := updated.Username

		if user.Username != *newUsername {
			taken, err := s.repository.IsUsernameTaken(*newUsername)
			if err != nil {
				return err
			}
			if !taken {
				user.Username = *newUsername
			} else {
				return errors.New(constants.UsernameIsTaken)
			}
		}
	}

	if updated.Name != nil {
		user.Name = *updated.Name
	}

	if updated.AvatarURL != nil {
		user.AvatarUrl = *updated.AvatarURL
	}

	return s.repository.Update(user)
}

func (s *service) UploadAvatarToS3(file *multipart.File, fileHeader *multipart.FileHeader) (*string, int, error) {
	maxSize_1MB := int64(1 * 1024 * 1024)
	if err := utils.IsValidImage(file, fileHeader, maxSize_1MB); err != nil {
		return nil, http.StatusBadRequest, err
	}

	avatarUrl, err := utils.UploadToS3(*file, fileHeader, "avatar")
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &avatarUrl, http.StatusOK, nil
}
