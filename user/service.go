package user

import (
	"avarts/constants"
	"avarts/utils"
	"errors"
	"mime/multipart"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type Service interface {
	GetUserByUsername(username string) (*User, error)
	GetUser(userID uint) (*User, error)
	DeleteUser(userID uint) (int, error)
	UpdateUser(userID uint, updated UpdateProfileRequest) (int, error)
	UploadAvatarToS3(file *multipart.File, fileHeader *multipart.FileHeader) (*string, int, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository}
}

func (s *service) GetUserByUsername(username string) (*User, error) {
	return s.repository.GetUserByUsername(username)
}

func (s *service) GetUser(userID uint) (*User, error) {
	return s.repository.GetUser(userID)
}

func (s *service) DeleteUser(userID uint) (int, error) {
	user, err := s.GetUser(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, errors.New(constants.UserNotFound)
		}
		return http.StatusInternalServerError, err
	}

	// Delete user's avatar from cloud storage
	avatarUrl := user.AvatarUrl
	key := strings.ReplaceAll(avatarUrl, constants.AwsS3PrefixUrl, "")
	utils.DeleteS3File(key)

	// Delete user's activity photos from cloud storage
	activityPhotoUrls, _ := s.repository.GetPictureURLsByUserID(userID)
	if len(activityPhotoUrls) > 0 {
		for _, url := range activityPhotoUrls {
			key := strings.ReplaceAll(url, constants.AwsS3PrefixUrl, "")
			utils.DeleteS3File(key)
		}
	}

	if err := s.repository.DeleteUser(userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, errors.New(constants.UserNotFound)
		}
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (s *service) UpdateUser(userID uint, updated UpdateProfileRequest) (int, error) {
	user, err := s.repository.GetUser(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, errors.New(constants.UserNotFound)
		}
		return http.StatusInternalServerError, err
	}

	if updated.Username != nil {
		newUsername := *updated.Username

		if err = utils.ValidateUsername(newUsername); err != nil {
			return http.StatusBadRequest, err
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

	if err = s.repository.UpdateUser(user); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (s *service) UploadAvatarToS3(file *multipart.File, fileHeader *multipart.FileHeader) (*string, int, error) {
	maxSize_300KB := int64(300 * 1024)
	if err := utils.IsValidImage(file, fileHeader, maxSize_300KB); err != nil {
		return nil, http.StatusBadRequest, err
	}

	avatarUrl, err := utils.UploadToS3(*file, fileHeader, "avatar")
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &avatarUrl, http.StatusOK, nil
}
