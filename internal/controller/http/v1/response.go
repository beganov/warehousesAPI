package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type response struct {
	Status        int    `json:"status"`
	StatusMessage string `json:"status_message"`
	Message       any    `json:"message"`
	Error         string `json:"error" example:"message"`
}

func errorResponse(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, response{Status: code,
		StatusMessage: http.StatusText(code),
		Error:         msg,
	})
}

func successResponse(c *gin.Context, code int, msg any) {
	c.JSON(code, response{
		Status:        code,
		StatusMessage: http.StatusText(code),
		Message:       msg,
	})
}
