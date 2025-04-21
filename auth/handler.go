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
	var req GoogleLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil || req.IdToken == "" {
		response.Failed(c, http.StatusBadRequest, constants.MissingIDToken)
		return
	}

	LoginResponse, err := h.service.LoginWithGoogle(req.IdToken)
	if err != nil {
		response.Failed(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.LoginSuccess, LoginResponse)
}
