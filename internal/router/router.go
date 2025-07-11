// Package router содержит маршруты gin
package router

import (
	"net/http"
	"vault-exporter/internal/config"
	"vault-exporter/internal/handler"
	"vault-exporter/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.ServerConfig) *gin.Engine {
	r := gin.Default()

	// Регистрируем все зависимости

	fileGetterService := service.NewFileGetterService(cfg)
	fgHandler := handler.NewLoadVaultHandler(fileGetterService, cfg)

	api := r.Group("/api")
	{
		fgHandler.RegisterRoute(api)
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Hello world"})
		})
	}

	return r
}
