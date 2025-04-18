package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	GoogleID  string    `gorm:"unique" json:"google_id"`
	Name      string    `json:"name"`
	Email     string    `gorm:"unique" json:"email"`
	AvatarUrl string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func MigrateUser(db *gorm.DB) {
	db.AutoMigrate(&User{})
}
