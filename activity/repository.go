package activity

import (
	"gorm.io/gorm"
)

type Repository interface {
	// Activity
	Create(activity *Activity) error
	GetByID(activityID *uint) (*Activity, error)
	GetAll(userID *uint) (*[]Activity, error)
	DeleteActivityByID(activityID uint) error

	// Picture
	CreatePicture(picture *Picture) error
	GetPictureUrlsByActivityID(activityID *uint) (*[]string, error)
	DeletePictureByID(id uint) error

	// Like
	IsLikeExists(like *Like) (bool, error)
	CreateLike(like *Like) error
	DeleteLike(like *Like) error

	// Comment
	GetCommentWithActivity(commentID uint) (*Comment, error)
	CreateComment(comment *Comment) error
	DeleteComment(commentID uint) error
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
	return r.db.Create(pics).Error
}

func (r *repository) GetByID(activityID *uint) (*Activity, error) {
	var activity Activity
	err := r.db.Preload("User").First(&activity, activityID).Error
	if err != nil {
		return nil, err
	}
	return &activity, nil
}

func (r *repository) GetAll(userID *uint) (*[]Activity, error) {
	var activities []Activity
	query := r.db.Preload("User")
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	err := query.Find(&activities).Error
	return &activities, err
}

func (r *repository) GetPictureUrlsByActivityID(activityID *uint) (*[]string, error) {
	var urls []string

	if err := r.db.
		Model(&Picture{}).
		Where("activity_id = ?", activityID).
		Pluck("url", &urls).Error; err != nil {
		return nil, err
	}

	return &urls, nil
}

func (r *repository) DeletePictureByID(id uint) error {
	return r.db.Delete(&Picture{}, id).Error
}

func (r *repository) DeleteActivityByID(activityID uint) error {
	result := r.db.Delete(&Activity{}, activityID)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (r *repository) CreateLike(like *Like) error {
	return r.db.Create(like).Error
}

func (r *repository) DeleteLike(like *Like) error {
	return r.db.Where("activity_id = ? AND user_id = ?", like.ActivityID, like.UserID).Delete(&Like{}).Error
}

func (r *repository) IsLikeExists(like *Like) (bool, error) {
	var count int64
	err := r.db.Model(&Like{}).Where("activity_id = ? AND user_id = ?", like.ActivityID, like.UserID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *repository) GetCommentWithActivity(commentID uint) (*Comment, error) {
	var comment Comment
	err := r.db.Preload("Activity").First(&comment, commentID).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *repository) CreateComment(comment *Comment) error {
	return r.db.Create(comment).Error
}

func (r *repository) DeleteComment(commentID uint) error {
	return r.db.Delete(&Comment{}, commentID).Error
}