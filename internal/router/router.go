// Package router содержит маршруты gin
package router

import (
	"vault-exporter/internal/config"
	"vault-exporter/internal/handler"
	"vault-exporter/internal/logger"
	"vault-exporter/internal/repository"
	"vault-exporter/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupServer(cfg *config.ServerConfig, db *pgxpool.Pool, logger logger.Logger) *gin.Engine {
	r := gin.New()
	r.Use(gin.LoggerWithWriter(logger.Writer()))
	r.Use(gin.Recovery())

	ksRepo := repository.NewKSRepository(db)

	// Регистрируем все зависимости

	fileGetterService := service.NewFileGetterService(cfg)
	izdCreatorService := service.NewIzdCreatorService(cfg, fileGetterService, ksRepo)
	importProcessorService := service.NewImportProcessorService(cfg, db, fileGetterService, izdCreatorService, ksRepo)

	fgHandler := handler.NewLoadVaultHandler(cfg, importProcessorService)

	api := r.Group("/api")
	{
		fgHandler.RegisterRoute(api)
	}

	return r
}
