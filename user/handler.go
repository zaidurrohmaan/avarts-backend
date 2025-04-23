package user

import (
	"avarts/constants"
	"avarts/response"
	"avarts/utils"
	"errors"
	"log"
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

	url, err := utils.UploadToS3(file, fileHeader, "avatar")
	if err != nil {
		log.Println("UploadToS3 error:", err) // Tambahkan log ini
		response.Failed(c, http.StatusInternalServerError, constants.FileUploadFailed)
		return
	}

	response.Success(c, http.StatusCreated, constants.FileUploadSuccess, gin.H{"photo_url": url})
}

func (h *Handler) Profile(c *gin.Context) {
	username := c.Param("username")

	user, err := h.service.GetByUsername(username)
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

func (h *Handler) MyProfile(c *gin.Context) {
	userID, isError := utils.GetUserIDFromJWT(c)
	if isError {
		return
	}

	user, err := h.service.GetByID(userID)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}
	userResponse := GenerateUserResponse(user)

	response.Success(c, http.StatusOK, constants.UserFetchSuccess, userResponse)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	userID, isError := utils.GetUserIDFromJWT(c)
	if isError {
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.service.UpdateProfile(userID, req)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.UserUpdateSuccess, nil)
}

// func (h *Handler) UploadAvatar(c *gin.Context) {
// 	file, fileHeader, err := c.Request.FormFile("avatar")
// 	if err != nil {
// 		response.Failed(c, http.StatusBadRequest, constants.PhotoFileRequired)
// 		return
// 	}

// 	url, err := utils.UploadToS3(file, fileHeader, "avatars")
// 	if err != nil {
// 		response.Failed(c, http.StatusInternalServerError, constants.FileUploadFailed)
// 		return
// 	}

// 	userID, isError := utils.GetUserIDFromJWT(c)
// 	if isError {
// 		return
// 	}

// 	updatedUser := &User {
// 		AvatarUrl: url,
// 	}

// 	if err := h.service.UpdateProfile(userID, updatedUser); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save avatar URL"})
// 		return
// 	}
// }
