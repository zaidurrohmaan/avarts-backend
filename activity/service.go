package activity

import "errors"

type Service interface {
	// Activity
	CreateActivity(activity *Activity) error
	GetByID(activityID *uint) (*Activity, error)
	GetAll(userID *uint) (*[]Activity, error)
	DeleteActivityByID(activityID uint) error

	// Picture
	CreatePicture(picture *Picture) error
	GetPictureUrlsByActivityID(activityID *uint) (*[]string, error)
	DeletePictureByID(id uint) error

	// Like
	CreateLike(like *Like) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository}
}

func (s *service) CreateActivity(activity *Activity) error {
	return s.repository.Create(activity)
}

func (s *service) CreatePicture(picture *Picture) error {
	return s.repository.CreatePicture(picture)
}

func (s *service) GetByID(activityID *uint) (*Activity, error) {
	return s.repository.GetByID(activityID)
}

func (s *service) GetAll(userID *uint) (*[]Activity, error) {
	return s.repository.GetAll(userID)
}

func (s *service) GetPictureUrlsByActivityID(activityID *uint) (*[]string, error) {
	return s.repository.GetPictureUrlsByActivityID(activityID)
}

func (s *service) DeletePictureByID(id uint) error {
	return s.repository.DeletePictureByID(id)
}

func (s *service) DeleteActivityByID(activityID uint) error {
	return s.repository.DeleteActivityByID(activityID)
}

func (s *service) CreateLike(like *Like) error {
	exist, err := s.repository.IsLikeExist(like.ActivityID, like.UserID)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("like already exists")
	}

	return s.repository.CreateLike(like)
}