package activity

import (
	"avarts/constants"
	"errors"
	"net/http"
)

type Service interface {
	// Activity
	CreateActivity(activity *Activity) error
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

func (s *service) CreateActivity(activity *Activity) error {
	return s.repository.Create(activity)
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

	return http.StatusForbidden, errors.New(constants.DELETE_COMMENT_ACCESS_DENIED)
}
