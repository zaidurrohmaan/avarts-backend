package repository

import (
	"avarts/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetByGoogleId(googleId string) (*models.User, error)
	Create(user *models.User) error
	GetById(id uint) (*models.User, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepo {
	return &userRepo{db}
}

func (r *userRepo) GetByGoogleId(googleId string) (*models.User, error) {
	var user models.User

	result := r.db.Where("google_id = ?", googleId).First(&user)
	return &user, result.Error
}

func (r *userRepo) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepo) GetById(id uint) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, id)
	return &user, result.Error
}