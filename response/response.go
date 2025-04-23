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
		Status:  status,
		Message: message,
		Data:    data,
	})
}

func Success(c *gin.Context, code int, message string, data interface{}) {
	SendResponse(c, code, constants.StatusSuccess, message, data)
}

func Failed(c *gin.Context, code int, message string) {
	SendResponse(c, code, constants.StatusFailed, message, nil)
}

func SendSuccessWithWarning(c *gin.Context, message string) {
	Success(c, http.StatusCreated, message, nil)
}
