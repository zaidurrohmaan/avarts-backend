package user

import (
	"avarts/constants"
	"avarts/response"
	"avarts/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

func (h *Handler) UploadAvatar(c *gin.Context) {
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

	url, statusCode, err := h.service.UploadAvatarToS3(&file, fileHeader)
	if err != nil {
		response.Failed(c, statusCode, err.Error())
		return
	}

	response.Success(c, statusCode, constants.FileUploadSuccess, gin.H{"photo_url": url})
}

func (h *Handler) GetUser(c *gin.Context) {
	username := c.Param("username")

	user, err := h.service.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Failed(c, http.StatusNotFound, constants.UserNotFound)
		} else {
			response.Failed(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	userResponse := GenerateUserResponse(user)

	response.Success(c, http.StatusOK, constants.UserFetchSuccess, userResponse)
}

func (h *Handler) GetMyProfile(c *gin.Context) {
	userID, isError := utils.GetUserIDFromJWT(c)
	if isError {
		return
	}

	user, err := h.service.GetUser(userID)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}
	userResponse := GenerateUserResponse(user)

	response.Success(c, http.StatusOK, constants.UserFetchSuccess, userResponse)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	userID, isError := utils.GetUserIDFromJWT(c)
	if isError {
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, http.StatusBadRequest, err.Error())
		return
	}

	statusCode, err := h.service.UpdateUser(userID, req)
	if err != nil {
		response.Failed(c, statusCode, err.Error())
		return
	}

	response.Success(c, statusCode, constants.UserUpdateSuccess, nil)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	userID, isError := utils.GetUserIDFromJWT(c)
	if isError {
		return
	}

	statusCode, err := h.service.DeleteUser(userID)
	if err != nil {
		response.Failed(c, statusCode, err.Error())
		return
	}
	response.Success(c, statusCode, constants.UserDeleteSuccess, nil)
}
