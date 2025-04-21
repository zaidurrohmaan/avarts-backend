package activity

import (
	"avarts/constants"
	"errors"
	"net/http"
	"time"
)

type Service interface {
	// Activity
	CreateActivity(userID uint, activity *CreateActivityRequest) (*uint, int, error)
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
	CreateComment(comment *Comment) (*CreateCommentResponse, error)
	DeleteComment(userID, commentID uint) (int, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository}
}

func (s *service) CreateActivity(userID uint, request *CreateActivityRequest) (*uint, int, error) {
	startTime, _ := time.Parse(time.RFC3339, request.StartTime)
	endTime, _ := time.Parse(time.RFC3339, request.EndTime)
	date, _ := time.Parse("2006-01-02", request.Date)

	activity := Activity{
		UserID:       userID,
		Title:        request.Title,
		Caption:      request.Caption,
		Distance:     request.Distance,
		Pace:         request.Pace,
		StepsCount:   request.StepsCount,
		StartTime:    startTime,
		EndTime:      endTime,
		Location:     request.Location,
		ActivityType: request.ActivityType,
		Date:         date,
	}

	// Save activity
	err := s.repository.Create(&activity)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Map PictureUrls to ActivityID
	for _, url := range request.PictureURLs {
		pic := Picture{
			ActivityID: activity.ID,
			URL:        url,
		}
		if err := s.repository.CreatePicture(&pic); err != nil {
			// Rollback: delete the activity, and all associated pictures will be automatically deleted as well
			_ = s.repository.DeleteActivityByID(activity.ID)
			return nil, http.StatusInternalServerError, errors.New(constants.ActivityCreateFailed)
		}
	}

	return &activity.ID, http.StatusCreated, nil
}

func (s *service) CreatePicture(picture *Picture) error {
	return s.repository.CreatePicture(picture)
}

func (s *service) GetByID(activityID *uint) (*Activity, error) {
	return s.repository.GetByID(activityID)
}

func (s *service) GetAll(userID *uint) (*[]Activity, error) {
	return s.repository.GetAll(userID)
}

func (s *service) GetPictureUrlsByActivityID(activityID *uint) (*[]string, error) {
	return s.repository.GetPictureUrlsByActivityID(activityID)
}

func (s *service) DeletePictureByID(id uint) error {
	return s.repository.DeletePictureByID(id)
}

func (s *service) DeleteActivityByID(activityID uint) error {
	return s.repository.DeleteActivityByID(activityID)
}

func (s *service) IsLikeExists(like *Like) (bool, error) {
	return s.repository.IsLikeExists(like)
}

func (s *service) CreateLike(like *Like) error {
	return s.repository.CreateLike(like)
}

func (s *service) DeleteLike(like *Like) error {
	return s.repository.DeleteLike(like)
}

func (s *service) CreateComment(comment *Comment) (*CreateCommentResponse, error) {
	err := s.repository.CreateComment(comment)
	if err != nil {
		return nil, err
	}

	response := &CreateCommentResponse{
		ID:         comment.ID,
		UserID:     comment.UserID,
		ActivityID: comment.ActivityID,
		Text:       comment.Text,
	}

	return response, nil
}

func (s *service) DeleteComment(userID, commentID uint) (int, error) {
	comment, err := s.repository.GetCommentWithActivity(commentID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	commentOwner := comment.UserID
	activityOwner := comment.Activity.UserID

	if userID == commentOwner || userID == activityOwner {
		err = s.repository.DeleteComment(commentID)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		return http.StatusOK, nil
	}

	return http.StatusForbidden, errors.New(constants.CommentDeleteAccessDenied)
}
