package auth

import (
	"avarts/utils"
	"errors"
	"fmt"
	"strings"

	userPackage "avarts/user"

	"gorm.io/gorm"
)

type Service interface {
	LoginWithGoogle(idToken string) (*LoginResponse, error)
}

type service struct {
	repository userPackage.Repository
}

func NewService(repository userPackage.Repository) Service {
	return &service{repository}
}

func (s *service) LoginWithGoogle(idToken string) (*LoginResponse, error) {
	googleUserInfo, err := utils.VerifyGoogleToken(idToken)
	if err != nil {
		return nil, err
	}

	var responseData LoginResponse

	user, err := s.repository.GetByGoogleId(googleUserInfo.GoogleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			emailPrefix := strings.SplitN(googleUserInfo.Email, "@", 2)
			baseUsername := strings.ToLower(emailPrefix[0])
			username := strings.ReplaceAll(baseUsername, "-", "")

			for i := 1; ; i++ {
				taken, err := s.repository.IsUsernameTaken(username)
				if err != nil {
					return nil, err
				}
				if !taken {
					break
				}
				username = fmt.Sprintf("%s%d", baseUsername, i)
			}

			user = &userPackage.User{
				Username:  username,
				Name:      googleUserInfo.Name,
				Email:     googleUserInfo.Email,
				AvatarUrl: googleUserInfo.Picture,
				GoogleID:  googleUserInfo.GoogleID,
			}

			err = s.repository.Create(user)
			if err != nil {
				return nil, err
			}

			token, err := utils.GenerateJWT(user.Username, user.ID)
			if err != nil {
				return nil, err
			}

			responseData.UserID = user.ID
			responseData.IsNew = true
			responseData.Token = token

			return &responseData, nil
		} else {
			return nil, err
		}
	}

	token, err := utils.GenerateJWT(user.Username, user.ID)
	if err != nil {
		return nil, err
	}

	responseData.UserID = user.ID
	responseData.IsNew = false
	responseData.Token = token

	return &responseData, nil
}
