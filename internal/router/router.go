// Package router содержит маршруты gin
package router

import (
	"net/http"
	"vault-exporter/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/load", handler.LoadVaultData)
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Hello world"})
		})
	}

	return r
}
