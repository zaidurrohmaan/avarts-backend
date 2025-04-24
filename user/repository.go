package user

import (
	"log"

	"gorm.io/gorm"
)

type Repository interface {
	GetUser(userID uint) (*User, error)
	GetUserByUsername(username string) (*User, error)
	GetUserByGoogleId(googleId string) (*User, error)
	CreateUser(user *User) error
	UpdateUser(user *User) error
	DeleteUser(userID uint) error
	IsUsernameTaken(username string) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) GetUser(userID uint) (*User, error) {
	var user User
	result := r.db.First(&user, userID)
	return &user, result.Error
}

func (r *repository) GetUserByUsername(username string) (*User, error) {
	var user User
	result := r.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		log.Println("Error in GetUserByUsername:", result.Error)
	}
	return &user, result.Error
}

func (r *repository) GetUserByGoogleId(googleId string) (*User, error) {
	var user User

	result := r.db.Where("google_id = ?", googleId).First(&user)
	return &user, result.Error
}

func (r *repository) CreateUser(user *User) error {
	return r.db.Create(user).Error
}

func (r *repository) UpdateUser(user *User) error {
	return r.db.Save(user).Error
}

func (r *repository) DeleteUser(userID uint) error {
	result := r.db.Delete(&User{}, userID)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *repository) IsUsernameTaken(username string) (bool, error) {
	var count int64
	err := r.db.Model(&User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}
