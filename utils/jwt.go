package utils

import (
	"avarts/config"
	"avarts/constants"
	"avarts/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(username string, userID uint) (string, error) {
	claims := jwt.MapClaims{
		"id":       userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWTSecret))
}

func ParseJWT(tokenStr string) (uint, string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})
	if err != nil {
		return 0, "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		id := uint(claims["id"].(float64))
		username := string(claims["username"].(string))
		return id, username, nil
	}
	return 0, "", err
}

func GetUserIDFromJWT(c *gin.Context) (uint, bool) {
	idInterface, exists := c.Get("id")
	if !exists {
		response.Failed(c, http.StatusUnauthorized, constants.Unauthorized)
		return 0, true
	}

	userId, ok := idInterface.(uint)
	if !ok {
		response.Failed(c, http.StatusInternalServerError, constants.InvalidRequestFormat)
		return 0, true
	}

	return userId, false
}
