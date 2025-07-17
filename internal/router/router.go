// Package router содержит маршруты gin
package router

import (
	"database/sql"
	"net/http"
	"vault-exporter/internal/config"
	"vault-exporter/internal/handler"
	"vault-exporter/internal/repository"
	"vault-exporter/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupServer(cfg *config.ServerConfig, db *sql.DB) *gin.Engine {
	r := gin.Default()

	ksRepo := repository.NewKSRepository(db)

	// Регистрируем все зависимости

	fileGetterService := service.NewFileGetterService(cfg)
	izdCreatorService := service.NewIzdCreatorService(cfg, &ksRepo)

	fgHandler := handler.NewLoadVaultHandler(fileGetterService, izdCreatorService, cfg)

	api := r.Group("/api")
	{
		fgHandler.RegisterRoute(api)
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Hello world"})
		})
	}

	return r
}
