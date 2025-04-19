package activity

import (
	"gorm.io/gorm"
)

type Repository interface {
	Create(activity *Activity) error
	CreatePicture(picture *Picture) error
	GetByID(id uint) (*Activity, error)
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

func (r *repository) GetByID(id uint) (*Activity, error) {
	var activity Activity
	result := r.db.First(&activity, id)
	return &activity, result.Error
}
