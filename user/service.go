package user

import (
	"fmt"
)

type Service interface {
	GetByUsername(username string) (*User, error)
	GetByID(userID uint) (*User, error)
	UpdateProfile(userID uint, updated *User) (*User, error)
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

func (s *service) UpdateProfile(userID uint, updated *User) (*User, error) {
	user, err := s.repository.Get(userID)
	if err != nil {
		return nil, err
	}

	newUsername := updated.Username

	if user.Username != newUsername {
		taken, err := s.repository.IsUsernameTaken(newUsername)
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

	err = s.repository.Update(user)
	return user, err
}