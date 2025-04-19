package activity

import (
	"gorm.io/gorm"
)

type Repository interface {
	Create(activity *Activity) error
	CreatePicture(picture *Picture) error
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
