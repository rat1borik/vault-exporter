// Package router содержит маршруты gin
package router

import (
	"net/http"
	"vault-exporter/internal/config"
	"vault-exporter/internal/handler"
	"vault-exporter/internal/repository"
	"vault-exporter/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupServer(cfg *config.ServerConfig, db *pgxpool.Pool) *gin.Engine {
	r := gin.Default()

	ksRepo := repository.NewKSRepository(db)

	// Регистрируем все зависимости

	fileGetterService := service.NewFileGetterService(cfg)
	izdCreatorService := service.NewIzdCreatorService(cfg, fileGetterService, ksRepo)
	importProcessorService := service.NewImportProcessorService(cfg, db, fileGetterService, izdCreatorService, ksRepo)

	fgHandler := handler.NewLoadVaultHandler(cfg, importProcessorService)

	api := r.Group("/api")
	{
		fgHandler.RegisterRoute(api)
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Hello world"})
		})
	}

	return r
}
