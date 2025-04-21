package activity

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
	IsLikeExists(like *Like) (bool, error)
	CreateLike(like *Like) error
	DeleteLike(like *Like) error
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

func (s *service) IsLikeExists(like *Like) (bool, error) {
	return s.repository.IsLikeExists(like)
}

func (s *service) CreateLike(like *Like) error {
	return s.repository.CreateLike(like)
}

func (s *service) DeleteLike(like *Like) error {
	return s.repository.DeleteLike(like)
}