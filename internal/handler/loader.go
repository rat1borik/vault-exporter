// Package handler содержит обработчики ручек.
package handler

import (
	"log"
	"vault-exporter/internal/config"
	"vault-exporter/internal/model"
	"vault-exporter/internal/response"
	"vault-exporter/internal/service"

	"github.com/gin-gonic/gin"
)

type LoadVaultHandler struct {
	fileGetterSvc service.FileGetterService
	cfg           *config.ServerConfig
}

func NewLoadVaultHandler(fileGetterSvc service.FileGetterService, cfg *config.ServerConfig) *LoadVaultHandler {
	return &LoadVaultHandler{
		fileGetterSvc: fileGetterSvc,
		cfg:           cfg,
	}
}

func (h *LoadVaultHandler) RegisterRoute(r *gin.RouterGroup) {
	r.POST("/load", h.LoadVaultData)
}

// Загружает данные при вызове ручки из Vault
func (h *LoadVaultHandler) LoadVaultData(c *gin.Context) {

	var items []model.VaultItem

	if err := c.ShouldBindBodyWithJSON(&items); err != nil {
		log.Printf("%v", err.Error())
		response.ValidationError(c, []string{"Ошибка при обработке входных данных"})
	}

	response.Success(c)
}
