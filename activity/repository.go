package activity

import (
	"gorm.io/gorm"
)

type Repository interface {
	Create(activity *Activity) error
	CreatePicture(picture *Picture) error
	GetByID(activityID uint) (*Activity, error)
	GetAll(userID *uint) (*[]Activity, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(activity *Activity) error {
	return r.db.Create(activity).Error
}

func (r *repository) CreatePicture(pics *Picture) error {
	return r.db.Create(&pics).Error
}

func (r *repository) GetByID(activityID uint) (*Activity, error) {
	var activity Activity
	result := r.db.First(&activity, activityID)
	return &activity, result.Error
}

func (r *repository) GetAll(userID *uint) (*[]Activity, error) {
	var activities []Activity
	query := r.db.Preload("User")
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	err := query.Find(&activities).Error
	return &activities, err
}
