package activity

import (
	"avarts/user"
	"time"

	"gorm.io/gorm"
)

type Activity struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint      `gorm:"not null" json:"user_id"`
	User         user.User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;" json:"user"`
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
	ID         uint     `gorm:"primaryKey;autoIncrement" json:"id"`
	ActivityID uint     `gorm:"not null" json:"activity_id"`
	URL        string   `gorm:"type:varchar(255)" json:"url"`
	Activity   Activity `gorm:"foreignKey:ActivityID;constraint:onDelete:CASCADE;" json:"-"`
}

type Like struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ActivityID uint      `gorm:"not null" json:"activity_id"`
	UserID     uint      `gorm:"not null" json:"user_id"`
	Activity   Activity  `gorm:"foreignKey:ActivityID;constraint:onDelete:CASCADE;" json:"-"`
	CreatedAt  time.Time `json:"created_at"`
}

type Comment struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ActivityID uint      `gorm:"not null" json:"activity_id"`
	UserID     uint      `gorm:"not null" json:"user_id"`
	Text       string    `gorm:"type:text;not null" json:"text"`
	Activity   Activity  `gorm:"foreignKey:ActivityID;constraint:OnDelete:CASCADE;" json:"-"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateCommentRequest struct {
	Text string `json:"text"`
}

type CreateCommentResponse struct {
	ID         uint   `json:"comment_id"`
	UserID     uint   `json:"user_id"`
	ActivityID uint   `json:"activity_id"`
	Text       string `json:"text"`
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

type CreateActivityResponse struct {
	ActivityID uint `json:"activity_id"`
}

type ActivityUserResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

type ActivityResponse struct {
	ID           uint                 `json:"id"`
	User         ActivityUserResponse `json:"user"`
	Title        string               `json:"title"`
	Caption      string               `json:"caption"`
	Distance     float64              `json:"distance"`
	Pace         float64              `json:"pace"`
	StepsCount   int                  `json:"steps_count"`
	StartTime    time.Time            `json:"start_time"`
	EndTime      time.Time            `json:"end_time"`
	Location     string               `json:"location"`
	ActivityType string               `json:"activity_type"`
	Date         time.Time            `json:"date"`
	PictureURLs  []string             `json:"picture_urls"`
}

func MigrateActivity(db *gorm.DB) {
	db.AutoMigrate(&Activity{}, &Picture{}, &Like{}, &Comment{})
}

func GenerateActivityResponse(activityData *Activity) ActivityResponse {
	response := ActivityResponse{
		ID: activityData.ID,
		User: ActivityUserResponse{
			ID:        activityData.User.ID,
			Username:  activityData.User.Username,
			Name:      activityData.User.Name,
			AvatarURL: activityData.User.AvatarUrl,
		},
		Title:        activityData.Title,
		Caption:      activityData.Caption,
		Distance:     activityData.Distance,
		Pace:         activityData.Pace,
		StepsCount:   activityData.StepsCount,
		StartTime:    activityData.StartTime,
		EndTime:      activityData.EndTime,
		Location:     activityData.Location,
		ActivityType: activityData.ActivityType,
		Date:         activityData.Date,
	}
	return response
}
