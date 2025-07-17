// Package handler содержит обработчики ручек.
package handler

import (
	"log"
	"vault-exporter/internal/config"
	"vault-exporter/internal/domain"
	"vault-exporter/internal/response"
	"vault-exporter/internal/service"

	"github.com/gin-gonic/gin"
)

type LoadVaultHandler struct {
	importProcSvc service.ImportProcessorService
	cfg           *config.ServerConfig
}

func NewLoadVaultHandler(cfg *config.ServerConfig, importProcSvc service.ImportProcessorService) *LoadVaultHandler {
	return &LoadVaultHandler{
		importProcSvc: importProcSvc,
		cfg:           cfg,
	}
}

func (h *LoadVaultHandler) RegisterRoute(r *gin.RouterGroup) {
	r.POST("/load", h.LoadVaultData)
}

// Загружает данные при вызове ручки из Vault
func (h *LoadVaultHandler) LoadVaultData(c *gin.Context) {

	var items []domain.VaultItem

	if err := c.ShouldBindBodyWithJSON(&items); err != nil {
		log.Printf("%v", err.Error())
		response.ValidationError(c, []string{"Ошибка при обработке входных данных"})
	}

	if err := h.importProcSvc.Import(items); err != nil {
		response.ServerError(c, []string{})
	}

	response.Success(c)
}
