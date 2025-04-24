package activity

import (
	"avarts/constants"
	"avarts/response"
	"avarts/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

func (h *Handler) UploadActivityPhoto(c *gin.Context) {
	fileHeader, err := c.FormFile("photo")
	if err != nil {
		response.Failed(c, http.StatusBadRequest, constants.PhotoFileRequired)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, constants.OpenFileFailed)
		return
	}
	defer file.Close()

	url, err := h.service.UploadActivityPhotoToS3(&file, fileHeader)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, constants.FileUploadFailed)
		return
	}

	response.Success(c, http.StatusCreated, constants.FileUploadSuccess, gin.H{"photo_url": url})
}

func (h *Handler) CreateActivity(c *gin.Context) {
	var req CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, http.StatusBadRequest, constants.InvalidRequestFormat)
		return
	}

	userID, isError := utils.GetUserIDFromJWT(c)
	if isError {
		return
	}

	activityID, statusCode, err := h.service.CreateActivity(userID, &req)
	if err != nil {
		response.Failed(c, statusCode, err.Error())
		return
	}

	responseData := &CreateActivityResponse{
		ActivityID: *activityID,
	}

	response.Success(c, http.StatusCreated, constants.ActivityCreateSuccess, responseData)
}

func (h *Handler) GetActivity(c *gin.Context) {
	activityID := utils.GetIDParam(c)
	if activityID == nil {
		return
	}

	responseData, statusCode, err := h.service.GetActivity(activityID)
	if err != nil {
		response.Failed(c, statusCode, err.Error())
		return
	}

	response.Success(c, statusCode, constants.ActivityFetchSuccess, responseData)
}

func (h *Handler) GetAllActivities(c *gin.Context) {
	var userID *uint
	if idStr := c.Query("userId"); idStr != "" {
		if id, err := strconv.ParseUint(idStr, 10, 64); err == nil {
			uid := uint(id)
			userID = &uid
		}
	}

	responseData, statusCode, err := h.service.GetAllActivities(userID)
	if err != nil {
		response.Failed(c, statusCode, err.Error())
		return
	}

	response.Success(c, statusCode, constants.ActivityFetchSuccess, responseData)
}

func (h *Handler) CreateLike(c *gin.Context) {
	userID, isError := utils.GetUserIDFromJWT(c)
	if isError {
		return
	}

	activityID := utils.GetIDParam(c)
	if activityID == nil {
		return
	}

	statusCode, err := h.service.CreateLike(userID, activityID)
	if err != nil {
		response.Failed(c, statusCode, err.Error())
		return
	}

	response.Success(c, statusCode, constants.LikeCreateSuccess, nil)
}

func (h *Handler) DeleteLike(c *gin.Context) {
	userID, isError := utils.GetUserIDFromJWT(c)
	if isError {
		return
	}

	activityID := utils.GetIDParam(c)
	if activityID == nil {
		return
	}

	statusCode, err := h.service.DeleteLike(userID, activityID)
	if err != nil {
		response.Failed(c, statusCode, err.Error())
		return
	}

	response.Success(c, statusCode, constants.LikeDeleted, nil)
}

func (h *Handler) CreateComment(c *gin.Context) {
	activityID := utils.GetIDParam(c)
	if activityID == nil {
		return
	}

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	userID, isError := utils.GetUserIDFromJWT(c)
	if isError {
		return
	}

	responseData, err := h.service.CreateComment(userID, activityID, &req)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, constants.CommentCreateSuccess, responseData)
}

func (h *Handler) DeleteComment(c *gin.Context) {
	userID, isError := utils.GetUserIDFromJWT(c)
	if isError {
		return
	}

	commentID := utils.GetIDParam(c)
	if commentID == nil {
		return
	}

	statusCode, err := h.service.DeleteComment(userID, *commentID)
	if err != nil {
		response.Failed(c, statusCode, err.Error())
		return
	}
	response.Success(c, statusCode, constants.CommentDeleteSuccess, nil)
}

func (h *Handler) DeleteActivity(c *gin.Context) {
	activityId := utils.GetIDParam(c)
	if activityId == nil {
		return
	}

	userID, isError := utils.GetUserIDFromJWT(c)
	if isError {
		return
	}

	statusCode, err := h.service.DeleteActivity(userID, activityId)
	if err != nil {
		response.Failed(c, statusCode, err.Error())
		return
	}

	response.Success(c, statusCode, constants.DeleteActivitySuccess, nil)
}
