package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type APIErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func Success(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, APIResponse{Success: true, Message: message, Data: data})
}

func Fail(c *gin.Context, status int, message string) {
	c.JSON(status, APIErrorResponse{Success: false, Message: message})
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
