package auth

import (
	"avarts/utils"
	"fmt"
	"strings"

	userPackage "avarts/user"
)

type Service interface {
	LoginWithGoogle(idToken, name, avatarUrl string) (string, error)
}

type service struct {
	repository userPackage.Repository
}

func NewService(repository userPackage.Repository) Service {
	return &service{repository}
}

func (s *service) LoginWithGoogle(idToken, name, avatarUrl string) (string, error) {
	tokenInfo, err := utils.VerifyGoogleToken(idToken)
	if err != nil {
		return "", err
	}

	user, err := s.repository.GetByGoogleId(tokenInfo.UserID)
	if err != nil {
		emailPrefix := strings.Split(tokenInfo.Email, "@")
		baseUsername := emailPrefix[0]
		username := baseUsername

		for i := 1; ; i++ {
			taken, err := s.repository.IsUsernameTaken(username)
			if err != nil {
				return "", err
			}
			if !taken {
				break
			}
			username = fmt.Sprintf("%s%d", baseUsername, i)
		}

		user = &userPackage.User{
			Username:  username,
			Name:      name,
			Email:     tokenInfo.Email,
			AvatarUrl: avatarUrl,
			GoogleID:  tokenInfo.UserID,
		}
		err = s.repository.Create(user)
		if err != nil {
			return "", err
		}
	}

	return utils.GenerateJWT(user.Username, user.ID)
}
