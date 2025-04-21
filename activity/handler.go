package activity

import (
	"avarts/constants"
	"avarts/response"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

func (h *Handler) UploadActivityPhoto(c *gin.Context) {
	file, err := c.FormFile("photo")
	if err != nil {
		response.SendError(c, http.StatusBadRequest, constants.PHOTO_FILE_REQUIRED)
		return
	}

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
	savePath := "./uploads/" + filename

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		response.SendError(c, http.StatusInternalServerError, constants.UPLOAD_FILE_FAILED)
		return
	}

	fileURL := fmt.Sprintf("http://localhost:8080/uploads/%s", filename)
	response.SendSuccess(c, http.StatusOK, constants.UPLOAD_FILE_SUCCESS, fileURL)
}

func (h *Handler) PostActivity(c *gin.Context) {
	var req CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, constants.INVALID_REQUEST)
		return
	}

	userIDInterface, _ := c.Get("id")
	userID := userIDInterface.(uint)

	startTime, _ := time.Parse(time.RFC3339, req.StartTime)
	endTime, _ := time.Parse(time.RFC3339, req.EndTime)
	date, _ := time.Parse("2006-01-02", req.Date)

	activity := Activity{
		UserID:       userID,
		Title:        req.Title,
		Caption:      req.Caption,
		Distance:     req.Distance,
		Pace:         req.Pace,
		StepsCount:   req.StepsCount,
		StartTime:    startTime,
		EndTime:      endTime,
		Location:     req.Location,
		ActivityType: req.ActivityType,
		Date:         date,
	}

	// Save activity
	if err := h.service.CreateActivity(&activity); err != nil {
		response.SendError(c, http.StatusInternalServerError, constants.CREATE_ACTIVITY_FAILED)
		return
	}

	// Map PictureUrls to ActivityID
	for _, url := range req.PictureURLs {
		pic := Picture{
			ActivityID: activity.ID,
			URL:        url,
		}
		if err := h.service.CreatePicture(&pic); err != nil {
			// Rollback: delete the activity, and all associated pictures will be automatically deleted as well
			_ = h.service.DeleteActivityByID(activity.ID)
			response.SendError(c, http.StatusInternalServerError, constants.CREATE_ACTIVITY_FAILED)
			return
		}
	}

	newActivity, err := h.service.GetByID(&activity.ID)
	if err != nil {
		response.SendSuccessWithWarning(c, constants.CREATE_ACTIVITY_SUCCESS_WITH_WARNING)
		return
	}
	newActivityResponse := GenerateActivityResponse(newActivity)

	pictureUrls, err := h.service.GetPictureUrlsByActivityID(&activity.ID)
	if err != nil {
		response.SendSuccessWithWarning(c, constants.CREATE_ACTIVITY_SUCCESS_WITH_WARNING)
		return
	}
	newActivityResponse.PictureURLs = *pictureUrls

	response.SendSuccess(c, http.StatusCreated, constants.CREATE_ACTIVITY_SUCCESS, newActivityResponse)
}

func (h *Handler) GetActivityByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.SendError(c, http.StatusBadRequest, constants.INVALID_TYPE_ACTIVITY_ID)
		return
	}

	activityID := uint(id)

	activity, err := h.service.GetByID(&activityID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.SendError(c, http.StatusNotFound, constants.ACTIVITY_NOT_FOUND)
			return
		} else {
			response.SendError(c, http.StatusInternalServerError, err.Error())
			return
		}
	}
	activityResponse := GenerateActivityResponse(activity)
	pictureUrls, err := h.service.GetPictureUrlsByActivityID(&activityID)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	activityResponse.PictureURLs = *pictureUrls

	response.SendSuccess(c, http.StatusOK, constants.ACTIVITY_FOUND_SUCCESS, activityResponse)
}

func (h *Handler) GetAllActivities(c *gin.Context) {
	var userID *uint
	if idStr := c.Query("userId"); idStr != "" {
		if id, err := strconv.ParseUint(idStr, 10, 64); err == nil {
			uid := uint(id)
			userID = &uid
		}
	}

	activities, err := h.service.GetAll(userID)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, constants.ACTIVITY_NOT_FOUND)
		return
	}

	var activitiesResponse []ActivityResponse
	for _, activity := range *activities {
		activityResponse := GenerateActivityResponse(&activity)
		activityID := activityResponse.ID
		pictureUrls, err := h.service.GetPictureUrlsByActivityID(&activityID)
		if err != nil {
			response.SendError(c, http.StatusInternalServerError, err.Error())
			return
		}
		activityResponse.PictureURLs = *pictureUrls
		activitiesResponse = append(activitiesResponse, activityResponse)
	}

	response.SendSuccess(c, http.StatusOK, constants.ACTIVITY_FOUND_SUCCESS, activitiesResponse)
}

func (h *Handler) CreateLike(c *gin.Context) {

	idInterface, exists := c.Get("id")
	if !exists {
		response.SendError(c, http.StatusUnauthorized, constants.UNAUTHORIZED)
		return
	}
	userId, ok := idInterface.(uint)
	if !ok {
		response.SendError(c, http.StatusInternalServerError, constants.INVALID_TYPE_USER_ID)
		return
	}

	var req LikeRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	like := &Like{
		ActivityID: req.ActivityID,
		UserID:     userId,
	}

	isLikeExists, err := h.service.IsLikeExists(like)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if isLikeExists {
		response.SendError(c, http.StatusBadRequest, constants.LIKE_ALREADY_EXISTS)
		return
	}

	err = h.service.CreateLike(like)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SendSuccess(c, http.StatusCreated, constants.CREATE_LIKE_SUCCESS, nil)
}

func (h *Handler) DeleteLike(c *gin.Context) {
	idInterface, exists := c.Get("id")
	if !exists {
		response.SendError(c, http.StatusUnauthorized, constants.UNAUTHORIZED)
		return
	}
	userId, ok := idInterface.(uint)
	if !ok {
		response.SendError(c, http.StatusInternalServerError, constants.INVALID_TYPE_USER_ID)
		return
	}

	var req LikeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	like := &Like{
		ActivityID: req.ActivityID,
		UserID: userId,
	}

	isLikeExists, err := h.service.IsLikeExists(like)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if !isLikeExists {
		response.SendError(c, http.StatusBadRequest, constants.LIKE_NOT_FOUND)
		return
	}

	err = h.service.DeleteLike(like)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SendSuccess(c, http.StatusOK, constants.LIKE_DELETED, nil)
}

func (h *Handler) CreateComment(c *gin.Context) {
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	idInterface, exists := c.Get("id")
	if !exists {
		response.SendError(c, http.StatusUnauthorized, constants.UNAUTHORIZED)
		return
	}
	userId, ok := idInterface.(uint)
	if !ok {
		response.SendError(c, http.StatusInternalServerError, constants.INVALID_TYPE_USER_ID)
		return
	}

	comment := &Comment {
		ActivityID: req.ActivityID,
		Text: req.Text,
		UserID: userId,
	}

	responseData, err := h.service.CreateComment(comment)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SendSuccess(c, http.StatusCreated, constants.CREATE_COMMENT_SUCCESS, responseData)
}

func (h *Handler) DeleteComment(c *gin.Context) {
	idInterface, exists := c.Get("id")
	if !exists {
		response.SendError(c, http.StatusUnauthorized, constants.UNAUTHORIZED)
		return
	}
	userId, ok := idInterface.(uint)
	if !ok {
		response.SendError(c, http.StatusInternalServerError, constants.INVALID_TYPE_USER_ID)
		return
	}

	var req DeleteCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	statusCode, err := h.service.DeleteComment(userId, req.CommentID)
	if err != nil {
		response.SendError(c, statusCode, err.Error())
		return
	}
	response.SendSuccess(c, statusCode, constants.DELETE_COMMENT_SUCCESS, nil)
}