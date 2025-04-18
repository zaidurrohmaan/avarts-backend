package controllers

import (
	"avarts/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service services.AuthService
}

func NewAuthController(service services.AuthService) *AuthController {
	return &AuthController{service}
}

type GoogleLoginRequest struct {
	IdToken    string `json:"id_token"`
	Name       string `json:"name"`
	AvatarURL  string `json:"avatar_url"`
}

func (ctrl *AuthController) GoogleLogin(c *gin.Context) {
	var body GoogleLoginRequest

	if err := c.ShouldBindJSON(&body); err != nil || body.IdToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing ID token"})
		return
	}

	token, err := ctrl.service.LoginWithGoogle(body.IdToken, body.Name, body.AvatarURL)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (ctrl *AuthController) Profile(c *gin.Context) {
	userId := c.MustGet("user_id").(uint)
	user, err := ctrl.service.GetProfile(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}