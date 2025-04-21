package utils

import (
	"avarts/constants"
	"avarts/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserIDFromJWT(c *gin.Context) (uint, bool) {
	idInterface, exists := c.Get("id")
	if !exists {
		response.SendError(c, http.StatusUnauthorized, constants.UNAUTHORIZED)
		return 0, true
	}

	userId, ok := idInterface.(uint)
	if !ok {
		response.SendError(c, http.StatusInternalServerError, constants.INVALID_TYPE_USER_ID)
		return 0, true
	}

	return userId, false
}