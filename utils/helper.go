package utils

import (
	"avarts/constants"
	"avarts/response"
	"errors"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ValidateUsername(username string) error {
	// Check length: must be between 1 and 30 characters
	if len(username) < 1 || len(username) > 30 {
		return errors.New("username must be between 1 and 30 characters")
	}

	// Allow only lower-case letters, numbers, periods, and underscores
	matched, _ := regexp.MatchString(`^[a-z0-9._]+$`, username)
	if !matched {
		return errors.New("username can only contain lower-case letters, numbers, dots, and underscores")
	}

	// Disallow consecutive periods or underscores (e.g., "..", "__", "._", "_.")
	if matched, _ := regexp.MatchString(`[._]{2,}`, username); matched {
		return errors.New("username cannot contain consecutive dots or underscores")
	}

	// Disallow starting or ending with a period or underscore
	if username[0] == '.' || username[0] == '_' || username[len(username)-1] == '.' || username[len(username)-1] == '_' {
		return errors.New("username cannot start or end with a dot or underscore")
	}

	return nil
}

func GetIDParam(c *gin.Context) (*uint) {
	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, constants.InvalidRequestFormat)
		return nil
	}

	id_uint := uint(id)
	return &id_uint
}