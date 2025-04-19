package activity

type Service interface {
	CreateActivity(activity *Activity) error
	CreatePicture(picture *Picture) error
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

