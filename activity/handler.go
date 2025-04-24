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

func (h *Handler) PostActivity(c *gin.Context) {
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

func (h *Handler) GetActivityByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, constants.InvalidRequestFormat)
		return
	}
	activityID := uint(id)

	responseData, statusCode, err := h.service.GetByID(&activityID)
	if err != nil {
		response.Failed(c, statusCode, err.Error())
		return
	}

	response.Success(c, statusCode, constants.ActivityFetchSuccess, responseData)
}

func (h *Handler) GetAllActivities(c *gin.Context) {
	var userID *uint
	if idStr := c.Query("userID"); idStr != "" {
		if id, err := strconv.ParseUint(idStr, 10, 64); err == nil {
			uid := uint(id)
			userID = &uid
		}
	}

	responseData, statusCode, err := h.service.GetAll(userID)
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

	var req LikeRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, err.Error())
		return
	}

	statusCode, err := h.service.CreateLike(userID, &req)
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

	var req LikeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	err := h.service.DeleteLike(userID, &req)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.LikeDeleted, nil)
}

func (h *Handler) CreateComment(c *gin.Context) {
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	userID, isError := utils.GetUserIDFromJWT(c)
	if isError {
		return
	}

	responseData, err := h.service.CreateComment(userID, &req)
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

	var req DeleteCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	statusCode, err := h.service.DeleteComment(userID, req.CommentID)
	if err != nil {
		response.Failed(c, statusCode, err.Error())
		return
	}
	response.Success(c, statusCode, constants.CommentDeleteSuccess, nil)
}

func (h *Handler) DeleteActivity(c *gin.Context) {
	userID, isError := utils.GetUserIDFromJWT(c)
	if isError {
		return
	}

	var req DeleteActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, http.StatusBadRequest, constants.InvalidRequestFormat)
		return
	}

	statusCode, err := h.service.DeleteActivity(userID, &req)
	if err != nil {
		response.Failed(c, statusCode, err.Error())
		return
	}

	response.Success(c, statusCode, constants.DeleteActivitySuccess, nil)
}
