package response

import (
	"avarts/constants"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func SendResponse(c *gin.Context, statusCode int, status string, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Status: status,
		Message: message,
		Data: data,
	})
}

func SendSuccess(c *gin.Context, code int, message string, data interface{}) {
	SendResponse(c, http.StatusOK, constants.STATUS_SUCCESS, message, data)
}

func SendError(c *gin.Context, code int, message string) {
	SendResponse(c, code, constants.STATUS_FAILED, message, nil)
}