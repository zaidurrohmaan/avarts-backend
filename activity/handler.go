package activity

import (
	"avarts/constants"
	"avarts/response"
	"avarts/utils"
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
		response.SendError(c, http.StatusBadRequest, constants.PhotoFileRequired)
		return
	}

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
	savePath := "./uploads/" + filename

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		response.SendError(c, http.StatusInternalServerError, constants.FileUploadFailed)
		return
	}

	fileURL := fmt.Sprintf("http://localhost:8080/uploads/%s", filename)
	response.SendSuccess(c, http.StatusOK, constants.FileUploadSuccess, fileURL)
}

func (h *Handler) PostActivity(c *gin.Context) {
	var req CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, constants.InvalidRequestFormat)
		return
	}

	userID, isError := utils.GetUserIDFromJWT(c)
	if isError {
		return
	}

	activityID, statusCode, err := h.service.CreateActivity(userID, &req)
	if err != nil {
		response.SendError(c, statusCode, err.Error())
	}

	responseData := &CreateActivityResponse {
		ActivityID: *activityID,
	}

	response.SendSuccess(c, http.StatusCreated, constants.ActivityCreateSuccess, responseData)
}

func (h *Handler) GetActivityByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.SendError(c, http.StatusBadRequest, constants.InvalidRequestFormat)
		return
	}

	activityID := uint(id)

	activity, err := h.service.GetByID(&activityID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.SendError(c, http.StatusNotFound, constants.ActivityNotFound)
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

	response.SendSuccess(c, http.StatusOK, constants.ActivityFetchSuccess, activityResponse)
}

func (h *Handler) GetAllActivities(c *gin.Context) {
	var userID *uint
	if idStr := c.Query("userID"); idStr != "" {
		if id, err := strconv.ParseUint(idStr, 10, 64); err == nil {
			uid := uint(id)
			userID = &uid
		}
	}

	activities, err := h.service.GetAll(userID)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, constants.ActivityNotFound)
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

	response.SendSuccess(c, http.StatusOK, constants.ActivityFetchSuccess, activitiesResponse)
}

func (h *Handler) CreateLike(c *gin.Context) {
	userID, isError := utils.GetUserIDFromJWT(c)
	if isError {
		return
	}

	var req LikeRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	statusCode, err := h.service.CreateLike(userID, &req)
	if err != nil {
		response.SendError(c, statusCode, err.Error())
	}

	response.SendSuccess(c, statusCode, constants.LikeCreateSuccess, nil)
}

func (h *Handler) DeleteLike(c *gin.Context) {
	userID, isError := utils.GetUserIDFromJWT(c)
	if isError {
		return
	}

	var req LikeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	err := h.service.DeleteLike(userID, &req)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SendSuccess(c, http.StatusOK, constants.LikeDeleted, nil)
}

func (h *Handler) CreateComment(c *gin.Context) {
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	userID, isError := utils.GetUserIDFromJWT(c)
	if isError {
		return
	}

	responseData, err := h.service.CreateComment(userID, &req)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SendSuccess(c, http.StatusCreated, constants.CommentCreateSuccess, responseData)
}

func (h *Handler) DeleteComment(c *gin.Context) {
	userID, isError := utils.GetUserIDFromJWT(c)
	if isError {
		return
	}

	var req DeleteCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	statusCode, err := h.service.DeleteComment(userID, req.CommentID)
	if err != nil {
		response.SendError(c, statusCode, err.Error())
		return
	}
	response.SendSuccess(c, statusCode, constants.CommentDeleteSuccess, nil)
}
