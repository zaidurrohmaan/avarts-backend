package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

func (handler *Handler) GoogleLogin(c *gin.Context) {
	var body GoogleLoginRequest

	if err := c.ShouldBindJSON(&body); err != nil || body.IdToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing ID token"})
		return
	}

	token, err := handler.service.LoginWithGoogle(body.IdToken, body.Name, body.AvatarURL)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}