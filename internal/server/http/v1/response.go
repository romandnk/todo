package v1

import (
	"github.com/gin-gonic/gin"
)

type response struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func sentErrorResponse(c *gin.Context, status int, message string, err error) {
	resp := response{
		Message: message,
		Error:   err.Error(),
	}

	c.AbortWithStatusJSON(status, resp)
}
