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

	var input User
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Failed(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.UpdateProfile(userID, &input)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, err.Error())
		return
	}
	userResponse := GenerateUserResponse(user)

	response.Success(c, http.StatusOK, constants.UserUpdateSuccess, userResponse)
}
