package services

import (
	"avarts/models"
	"avarts/repository"
	"avarts/utils"
)

type AuthService interface {
	LoginWithGoogle(idToken, name, avatarUrl string) (string, error)
	GetProfile(userId uint) (*models.User, error)
}

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) *authService {
	return &authService{repo}
}

func (s *authService) LoginWithGoogle(idToken, name, avatarUrl string) (string, error) {
	tokenInfo, err := utils.VerifyGoogleToken(idToken)
	if err != nil {
		return "", err
	}

	user, err := s.repo.GetByGoogleId(tokenInfo.UserID)
	if err != nil {
		user = &models.User{
			GoogleID: tokenInfo.UserID,
			Email:    tokenInfo.Email,
			Name: name,
			AvatarUrl: avatarUrl,
		}
		err = s.repo.Create(user)
		if err != nil {
			return "", err
		}
	}

	return utils.GenerateJWT(user.ID)
}

func (s *authService) GetProfile(userId uint) (*models.User, error) {
	return s.repo.GetById(userId)
}
