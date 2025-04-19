package activity

type Service interface {
	CreateActivity(activity *Activity) error
	CreatePicture(picture *Picture) error
	GetByID(activityID uint) (*Activity, error)
	GetAll(userID *uint) (*[]Activity, error)
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

func (s *service) GetByID(activityID uint) (*Activity, error) {
	return s.repository.GetByID(activityID)
}

func (s *service) GetAll(userID *uint) (*[]Activity, error) {
	return s.repository.GetAll(userID)
}