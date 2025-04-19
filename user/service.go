package user

import (
	"fmt"
)

type Service interface {
	GetProfile(username string) (*User, error)
	MyProfile(userId uint) (*User, error)
	UpdateProfile(userId uint, updated *User) (*User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository}
}

func (s *service) GetProfile(username string) (*User, error) {
	return s.repository.GetByUsername(username)
}

func (s *service) MyProfile(userId uint) (*User, error) {
	return s.repository.Get(userId)
}

func (s *service) UpdateProfile(userId uint, updated *User) (*User, error) {
	user, err := s.repository.Get(userId)
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