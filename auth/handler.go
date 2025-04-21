package auth

import (
	"avarts/constants"
	"avarts/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

func (h *Handler) GoogleLogin(c *gin.Context) {
	var body GoogleLoginRequest

	if err := c.ShouldBindJSON(&body); err != nil || body.IdToken == "" {
		response.SendError(c, http.StatusBadRequest, constants.MissingIDToken)
		return
	}

	token, err := h.service.LoginWithGoogle(body.IdToken, body.Name, body.AvatarURL)
	if err != nil {
		response.SendError(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.SendSuccess(c, http.StatusOK, constants.LoginSuccess, token)
}