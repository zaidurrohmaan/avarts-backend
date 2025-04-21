package utils

import (
	"avarts/constants"
	"avarts/response"
	"net/http"
	"os"
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
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ParseJWT(tokenStr string) (uint, string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
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
		response.SendError(c, http.StatusUnauthorized, constants.Unauthorized)
		return 0, true
	}

	userId, ok := idInterface.(uint)
	if !ok {
		response.SendError(c, http.StatusInternalServerError, constants.InvalidRequestFormat)
		return 0, true
	}

	return userId, false
}
