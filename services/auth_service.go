package services

import (
	"avarts/models"
	"avarts/repository"
	"avarts/utils"
	"fmt"
	"strings"
)

type AuthService interface {
	LoginWithGoogle(idToken, name, avatarUrl string) (string, error)
	GetProfile(username string) (*models.User, error)
	UpdateProfile(userId uint, updated *models.User) (*models.User, error)
}

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{repo}
}

func (s *authService) LoginWithGoogle(idToken, name, avatarUrl string) (string, error) {
	tokenInfo, err := utils.VerifyGoogleToken(idToken)
	if err != nil {
		return "", err
	}

	user, err := s.repo.GetByGoogleId(tokenInfo.UserID)
	if err != nil {
		emailPrefix := strings.Split(tokenInfo.Email, "@")
		baseUsername := emailPrefix[0]
		username := baseUsername

		for i := 1; ; i++ {
			taken, err := s.repo.IsUsernameTaken(username) 
			if err != nil {
				return "", err
			}
			if !taken {
				break
			}
			username = fmt.Sprintf("%s%d", baseUsername, i)
		}

		user = &models.User{
			Username:  username,
			Name:      name,
			Email:     tokenInfo.Email,
			AvatarUrl: avatarUrl,
			GoogleID:  tokenInfo.UserID,
		}
		err = s.repo.Create(user)
		if err != nil {
			return "", err
		}
	}

	return utils.GenerateJWT(user.Username, user.ID)
}

func (s *authService) GetProfile(username string) (*models.User, error) {
	return s.repo.GetByUsername(username)
}

func (s *authService) UpdateProfile(userId uint, updated *models.User) (*models.User, error) {
	user, err := s.repo.Get(userId)
	if err != nil {
		return nil, err
	}

	newUsername := updated.Username

	if user.Username != newUsername {
		taken, err := s.repo.IsUsernameTaken(newUsername)
		if err != nil {
			return nil, err
		}
		if !taken {
			user.Username = newUsername
		} else {
			return nil, fmt.Errorf("username already taken")
		}
	}

	user.Name = updated.Name
	user.AvatarUrl = updated.AvatarUrl

	err = s.repo.Update(user)
	return user, err
}
