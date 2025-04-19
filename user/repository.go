package user

import (
	"log"

	"gorm.io/gorm"
)

type Repository interface {
	Get(userId uint) (*User, error)
	GetByUsername(username string) (*User, error)
	GetByGoogleId(googleId string) (*User, error)
	Create(user *User) error
	Update(user *User) error
	IsUsernameTaken(username string) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Get(userId uint) (*User, error) {
	var user User
	result := r.db.First(&user, userId)
	return &user, result.Error
}

func (r *repository) GetByUsername(username string) (*User, error) {
	var user User
	result := r.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		log.Println("Error in GetByUsername:", result.Error)
	}
	return &user, result.Error
}

func (r *repository) GetByGoogleId(googleId string) (*User, error) {
	var user User

	result := r.db.Where("google_id = ?", googleId).First(&user)
	return &user, result.Error
}

func (r *repository) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *repository) Update(user *User) error {
	return r.db.Save(user).Error
}

func (r *repository) IsUsernameTaken(username string) (bool, error) {
	var count int64
	err := r.db.Model(&User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}
