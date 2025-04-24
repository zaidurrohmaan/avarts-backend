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
	UpdateProfile(userID uint, updated UpdateProfileRequest) (int, error)
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

func (s *service) UpdateProfile(userID uint, updated UpdateProfileRequest) (int, error) {
	user, err := s.repository.Get(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, errors.New(constants.UserNotFound)
		}
		return http.StatusInternalServerError, err
	}

	if updated.Username != nil {
		newUsername := *updated.Username

		if err = utils.ValidateUsername(newUsername); err != nil {
			return http.StatusBadRequest,err
		}

		if user.Username != newUsername {
			taken, err := s.repository.IsUsernameTaken(newUsername)
			if err != nil {
				return http.StatusInternalServerError, err
			}
			if !taken {
				user.Username = newUsername
			} else {
				return http.StatusBadRequest, errors.New(constants.UsernameIsTaken)
			}
		}
	}

	if updated.Name != nil {
		user.Name = *updated.Name
	}

	if updated.AvatarURL != nil {
		user.AvatarUrl = *updated.AvatarURL
	}

	if err = s.repository.Update(user); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
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
