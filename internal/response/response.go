// Package response содержит шорткаты для ответов сервера волту
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ServerError(c *gin.Context, messages []string) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"errors":  messages,
	})
}

func ValidationError(c *gin.Context, messages []string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"errors":  messages,
	})
}

func Success(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
