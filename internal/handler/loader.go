// Package handler содержит обработчики ручек.
package handler

import (
	"context"
	"log"
	"vault-exporter/internal/config"
	"vault-exporter/internal/domain"
	"vault-exporter/internal/response"
	"vault-exporter/internal/service"
	"vault-exporter/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		return
	}

	procId := uuid.New()

	procCtx := context.WithValue(c, utils.CtxProcId, procId.String())

	if err := h.importProcSvc.Import(procCtx, items); err != nil {
		msgs := make([]string, 0, len(err))
		for i := range err {
			msgs = append(msgs, err[i].Error())
		}
		response.ServerError(c, msgs)
		return
	}

	response.Success(c)
}
