package activity

import (
	"avarts/constants"
	"avarts/utils"
	"errors"
	"mime/multipart"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type Service interface {
	// Activity
	CreateActivity(userID uint, activity *CreateActivityRequest) (*uint, int, error)
	GetActivity(activityID *uint) (*ActivityResponse, int, error)
	GetAllActivities(userID *uint) (*[]ActivityResponse, int, error)
	UploadActivityPhotoToS3(file *multipart.File, fileHeader *multipart.FileHeader) (*string, error)
	DeleteActivity(userID uint, request *DeleteActivityRequest) (int, error)

	// Like
	CreateLike(userID uint, like *LikeRequest) (int, error)
	DeleteLike(userID uint, like *LikeRequest) (int, error)

	// Comment
	CreateComment(userID uint, request *CreateCommentRequest) (*CreateCommentResponse, error)
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
	err := s.repository.CreateActivity(&activity)
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
			_ = s.repository.DeleteActivity(activity.ID)
			return nil, http.StatusInternalServerError, errors.New(constants.ActivityCreateFailed)
		}
	}

	return &activity.ID, http.StatusCreated, nil
}

func (s *service) GetActivity(activityID *uint) (*ActivityResponse, int, error) {
	activity, err := s.repository.GetActivity(activityID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusNotFound, errors.New(constants.ActivityNotFound)
		} else {
			return nil, http.StatusInternalServerError, err
		}
	}
	activityResponse := GenerateActivityResponse(activity)
	pictureUrls, err := s.repository.GetPictureUrlsByActivityID(activityID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	activityResponse.PictureURLs = *pictureUrls

	return &activityResponse, http.StatusOK, nil
}

func (s *service) GetAllActivities(userID *uint) (*[]ActivityResponse, int, error) {
	activities, err := s.repository.GetAllActivities(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusNotFound, errors.New(constants.ActivityNotFound)
		}
		return nil, http.StatusInternalServerError, err
	}

	var activitiesResponse []ActivityResponse
	for _, activity := range *activities {
		activityResponse := GenerateActivityResponse(&activity)
		activityID := activityResponse.ID
		pictureUrls, err := s.repository.GetPictureUrlsByActivityID(&activityID)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		activityResponse.PictureURLs = *pictureUrls
		activitiesResponse = append(activitiesResponse, activityResponse)
	}

	return &activitiesResponse, http.StatusOK, nil
}

func (s *service) CreateLike(userID uint, request *LikeRequest) (int, error) {
	like := &Like{
		ActivityID: request.ActivityID,
		UserID:     userID,
	}

	isLikeExists, err := s.repository.IsLikeExists(like)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if isLikeExists {
		return http.StatusBadRequest, errors.New(constants.LikeAlreadyExists)
	}

	err = s.repository.CreateLike(like)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusCreated, nil
}

func (s *service) DeleteLike(userID uint, request *LikeRequest) (int, error) {
	like := &Like{
		ActivityID: request.ActivityID,
		UserID:     userID,
	}

	isLikeExists, err := s.repository.IsLikeExists(like)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if !isLikeExists {
		return http.StatusNotFound, errors.New(constants.LikeNotFound)
	}

	err = s.repository.DeleteLike(like)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (s *service) CreateComment(userID uint, request *CreateCommentRequest) (*CreateCommentResponse, error) {
	comment := &Comment{
		ActivityID: request.ActivityID,
		Text:       request.Text,
		UserID:     userID,
	}

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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, errors.New(constants.CommentNotFound)
		}
		return http.StatusInternalServerError, err
	}

	commentOwner := comment.UserID
	activityOwner := comment.Activity.UserID

	if userID == commentOwner || userID == activityOwner {
		err = s.repository.DeleteComment(commentID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return http.StatusNotFound, errors.New(constants.CommentNotFound)
			}
			return http.StatusInternalServerError, err
		}
		return http.StatusOK, nil
	}

	return http.StatusForbidden, errors.New(constants.CommentDeleteAccessDenied)
}

func (s *service) UploadActivityPhotoToS3(file *multipart.File, fileHeader *multipart.FileHeader) (*string, error) {
	maxSize_5MB := int64(5 * 1024 * 1024)
	if err := utils.IsValidImage(file, fileHeader, maxSize_5MB); err != nil {
		return nil, err
	}

	avatarUrl, err := utils.UploadToS3(*file, fileHeader, "activity")
	if err != nil {
		return nil, err
	}

	return &avatarUrl, nil
}

func (s *service) DeleteActivity(userID uint, request *DeleteActivityRequest) (int, error) {
	activity, statusCode, err := s.GetActivity(&request.ActivityID)
	if err != nil {
		return statusCode, err
	}

	if activity.User.ID != userID {
		return http.StatusForbidden, errors.New(constants.ActivityDeleteAccessDenied)
	}

	if err := s.repository.DeleteActivity(request.ActivityID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, errors.New(constants.ActivityNotFound)
		}
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
