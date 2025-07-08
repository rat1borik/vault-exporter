package router

import (
	"vault-exporter/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/load", handler.LoadVaultData)
	}

	return r
}
