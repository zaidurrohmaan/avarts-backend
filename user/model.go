package user

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"unique" json:"username"`
	Name      string    `json:"name"`
	Email     string    `gorm:"unique" json:"email"`
	AvatarUrl string    `json:"avatar_url"`
	GoogleID  string    `gorm:"unique" json:"google_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	AvatarUrl string    `json:"avatar_url"`
}

type UpdateProfileRequest struct {
	Username  *string `json:"username,omitempty"`
	Name      *string `json:"name,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

type DeleteUserRequest struct {
	UserID uint `json:"user_id"`
}

func GenerateUserResponse(userData *User) (UserResponse) {
	userResponse := UserResponse {
		ID: userData.ID,
		Username: userData.Username,
		Name: userData.Name,
		Email: userData.Email,
		AvatarUrl: userData.AvatarUrl,
	}
	return userResponse
}

func MigrateUser(db *gorm.DB) {
	db.AutoMigrate(&User{})
}