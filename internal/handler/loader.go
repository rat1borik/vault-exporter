// Package handler содержит обработчики ручек.
package handler

import (
	"bytes"
	"context"
	"io"
	"log"
	"sync"
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
	mu            sync.Mutex
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
	// Сериализуем запросы, чтобы поддерживать консистентность файлов в КС-ной папке (пусть ждут)
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.cfg.Server.ApiKey != "" {
		if token, err := getBearerToken(c); err != nil {
			log.Printf("%v", err.Error())
			response.Error(c, []string{"Ошибка получения ключа авторизации (запрос неверен)"})
			return
		} else if token != h.cfg.Server.ApiKey {
			response.Error(c, []string{"Ключ авторизации неверен"})
			return
		}
	}

	var items []domain.VaultItem

	if !h.cfg.IsProduction {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(400, gin.H{"error": "can't read body"})
			return
		}

		// Логируем тело
		log.Printf("Request body: %s\n", string(bodyBytes))

		// Восстанавливаем тело обратно в c.Request.Body
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	if err := c.ShouldBindBodyWithJSON(&items); err != nil {
		log.Printf("%v", err.Error())
		response.Error(c, []string{"Ошибка при обработке входных данных"})
		return
	}

	procId := uuid.New()

	procCtx := context.WithValue(c, utils.CtxProcId, procId.String())

	if err := h.importProcSvc.Import(procCtx, items); err != nil {
		msgs := make([]string, 0, len(err))
		for i := range err {
			log.Printf("%v", err[i].Error())
			msgs = append(msgs, err[i].Error())
		}
		response.Error(c, msgs)
		return
	}

	response.Success(c)
}
