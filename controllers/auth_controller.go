package controllers

import (
	"avarts/models"
	"avarts/services"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	username := c.Param("username")
	fmt.Println("Username from URL param:", username)
	user, err := ctrl.service.GetProfile(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, user)
}

func (ctrl *AuthController) UpdateProfile(c *gin.Context) {
	idInterface, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, ok := idInterface.(uint)
	if !ok {
    	c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
    	return
	}

	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctrl.service.UpdateProfile(userID, &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}