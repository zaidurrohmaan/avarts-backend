package user

import (
	"avarts/constants"
	"avarts/response"
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

	user, err := h.service.GetProfile(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.SendError(c, http.StatusNotFound, constants.USER_NOT_FOUND)
		} else {
			response.SendError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response.SendSuccess(c, http.StatusOK, constants.USER_FOUND_SUCCESS, user)
}

func (h *Handler) MyProfile(c *gin.Context) {
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

	user, err := h.service.MyProfile(userId)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SendSuccess(c, http.StatusOK, constants.USER_FOUND_SUCCESS, user)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	idInterface, exists := c.Get("id")
	if !exists {
		response.SendError(c, http.StatusUnauthorized, constants.UNAUTHORIZED)
		return
	}
	userID, ok := idInterface.(uint)
	if !ok {
		response.SendError(c, http.StatusInternalServerError, constants.INVALID_TYPE_USER_ID)
		return
	}

	var input User
	if err := c.ShouldBindJSON(&input); err != nil {
		response.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.UpdateProfile(userID, &input)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SendSuccess(c, http.StatusOK, constants.USER_UPDATE_SUCCESS, user)
}