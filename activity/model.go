package activity

import (
	"avarts/user"
	"time"

	"gorm.io/gorm"
)

type Activity struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint      `gorm:"not null" json:"user_id"`
	User         user.User `gorm:"foreignKey:UserID;references:ID" json:"user"`
	Title        string    `gorm:"type:varchar(100)" json:"title"`
	Caption      string    `gorm:"type:text" json:"caption"`
	Distance     float64   `gorm:"type:float" json:"distance"`
	Pace         float64   `gorm:"type:float" json:"pace"`
	StepsCount   int       `gorm:"type:int" json:"steps_count"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Location     string    `gorm:"type:varchar(255)" json:"location"`
	ActivityType string    `gorm:"type:varchar(50)" json:"activity_type"`
	Date         time.Time `gorm:"type:date" json:"date"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Picture struct {
	ID         uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	ActivityID uint   `gorm:"not null" json:"activity_id"`
	URL        string `gorm:"type:varchar(255)" json:"url"`
}

type CreateActivityRequest struct {
	Title        string   `json:"title"`
	Caption      string   `json:"caption"`
	Distance     float64  `json:"distance"`
	Pace         float64  `json:"pace"`
	StepsCount   int      `json:"steps_count"`
	StartTime    string   `json:"start_time"`
	EndTime      string   `json:"end_time"`
	Location     string   `json:"location"`
	ActivityType string   `json:"activity_type"`
	Date         string   `json:"date"`
	PictureURLs  []string `json:"picture_urls"`
}

func MigrateActivity(db *gorm.DB) {
	db.AutoMigrate(&Activity{}, &Picture{})
}
