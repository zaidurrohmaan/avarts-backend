package repository

import (
	"avarts/models"
	"log"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetByGoogleId(googleId string) (*models.User, error)
	Create(user *models.User) error
	Get(userId uint) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	Update(user *models.User) error
	IsUsernameTaken(username string) (bool, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) GetByGoogleId(googleId string) (*models.User, error) {
	var user models.User

	result := r.db.Where("google_id = ?", googleId).First(&user)
	return &user, result.Error
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Get(userId uint) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, userId)
	return &user, result.Error
}

func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	result := r.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		log.Println("Error in GetByUsername:", result.Error)
	}
	return &user, result.Error
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) IsUsernameTaken(username string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}
