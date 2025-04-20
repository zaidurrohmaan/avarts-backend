package activity

import (
	"gorm.io/gorm"
)

type Repository interface {
	Create(activity *Activity) error
	CreatePicture(picture *Picture) error
	GetByID(activityID *uint) (*Activity, error)
	GetAll(userID *uint) (*[]Activity, error)
	GetPictureUrlsByActivityID(activityID *uint) (*[]string, error)
	DeletePictureByID(id uint) error
	DeleteActivityByID(activityID uint) error
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

func (r *repository) GetByID(activityID *uint) (*Activity, error) {
	var activity Activity
	err := r.db.Preload("User").First(&activity, activityID).Error
	if err != nil {
		return nil, err
	}
	return &activity, nil
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

func (r *repository) GetPictureUrlsByActivityID(activityID *uint) (*[]string, error) {
	var urls []string

	if err := r.db.
		Model(&Picture{}).
		Where("activity_id = ?", activityID).
		Pluck("url", &urls).Error; err != nil {
		return nil, err
	}

	return &urls, nil
}

func (r *repository) DeletePictureByID(id uint) error {
	return r.db.Delete(&Picture{}, id).Error
}

func (r *repository) DeleteActivityByID(activityID uint) error {
	return r.db.Delete(&Activity{}, activityID).Error
}
