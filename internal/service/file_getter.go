// Package service содержит слой сервисов
package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"vault-exporter/internal/config"
	"vault-exporter/internal/model"
	"vault-exporter/internal/utils"
)

type FileGetterService interface {
	LoadFile(ctx context.Context, file *model.VaultFile) (string, error)
}

type fileGetterService struct {
	cfg *config.ServerConfig
}

func NewFileGetterService(cfg *config.ServerConfig) FileGetterService {
	return &fileGetterService{cfg: cfg}
}

// Осуществляет загрузку файла из Vault в специальную директорию КС
func (srv *fileGetterService) LoadFile(ctx context.Context, file *model.VaultFile) (string, error) {

	path := fmt.Sprintf("http://%s:%d/api/files?id=%d", srv.cfg.Vault.Host, srv.cfg.Vault.Port, file.Id)

	resp, err := http.Get(path)
	if err != nil {
		return "", fmt.Errorf("сan't perform request to vault: %w, id = %d", err, file.Id)
	}
	defer resp.Body.Close()

	fileName := file.FileName

	// Добиваемся уникальности
	if utils.FileExists(filepath.Join(srv.cfg.KSFilesPath, fileName)) {
		fileName = fmt.Sprintf("%s - %d", fileName, file.VerNum)

		if utils.FileExists(filepath.Join(srv.cfg.KSFilesPath, fileName)) {
			fileName = fmt.Sprintf("%s - %s", file.FileName, time.Now().Format("02.01.2006 15:04"))
		}
	}

	// Пишем во временное место (потом заберем)
	dest, err := os.Create(filepath.Join(srv.cfg.TempPath, fileName))
	if err != nil {
		return "", fmt.Errorf("сan't create file: %w, id = %d", err, file.Id)
	}

	_, err = io.Copy(dest, resp.Body)
	if err != nil {
		return "", fmt.Errorf("сan't write file: %w, id = %d", err, file.Id)
	}

	return fileName, nil
}
