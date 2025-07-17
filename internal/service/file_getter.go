// Package service содержит слой сервисов
package service

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"vault-exporter/internal/config"
	"vault-exporter/internal/model"
	"vault-exporter/internal/utils"
)

type FileGetterService interface {
	LoadFile(*model.VaultFile, string) (string, error)
	ClearTempFolder(string) error
	CommitFiles(string) error
	EnsureUnique(string) bool
}

type fileGetterService struct {
	cfg      *config.ServerConfig
	commitMu *sync.Mutex
}

func NewFileGetterService(cfg *config.ServerConfig) FileGetterService {
	return &fileGetterService{cfg: cfg}
}

// Осуществляет загрузку файла из Vault в специальную директорию КС
func (srv *fileGetterService) LoadFile(file *model.VaultFile, ctxId string) (string, error) {

	path := fmt.Sprintf("http://%s:%d/api/files?id=%d", srv.cfg.Vault.Host, srv.cfg.Vault.Port, file.Id)

	resp, err := http.Get(path)
	if err != nil {
		return "", fmt.Errorf("сan't perform request to vault: %w, id = %d", err, file.Id)
	}
	defer resp.Body.Close()

	// Пишем во временное место (потом заберем)
	utils.EnsureDir(srv.tempDirPath(ctxId))
	dest, err := os.Create(filepath.Join(srv.tempDirPath(ctxId), file.FileName))
	if err != nil {
		return "", fmt.Errorf("сan't create file: %w, id = %d", err, file.Id)
	}

	_, err = io.Copy(dest, resp.Body)
	if err != nil {
		return "", fmt.Errorf("сan't write file: %w, id = %d", err, file.Id)
	}

	return file.FileName, nil
}

func (srv *fileGetterService) tempDirPath(ctxId string) string {
	return filepath.Join(srv.cfg.TempPath, ctxId)
}

// Чистит временнную папку в случае неудачи / после всех действий
func (srv *fileGetterService) ClearTempFolder(ctxId string) error {
	err := utils.ClearDir(srv.tempDirPath(ctxId), true)
	if err != nil {
		return fmt.Errorf("can't clear temp dir: %w", err)
	}

	return nil
}

// Заливает файлы в финальную папку КС
func (srv *fileGetterService) CommitFiles(ctxId string) error {
	// Блокируемся, чтобы не было коллизий между двумя коммитами
	srv.commitMu.Lock()
	defer srv.commitMu.Unlock()

	entries, err := os.ReadDir(srv.tempDirPath(ctxId))
	if err != nil {
		return fmt.Errorf("can't commit files: %w", err)
	}

	for _, e := range entries {
		//TODO: проверка на уникальность

		utils.CopyFile(filepath.Join(srv.tempDirPath(ctxId), e.Name()), filepath.Join(srv.cfg.KSFilesPath, e.Name()), false)
	}

	srv.ClearTempFolder(ctxId)

	return nil
}

// Проверяет нет ли коллизий по данному ctxId (должен вызываться в самом начале )
func (srv *fileGetterService) EnsureUnique(ctxId string) bool {
	return !utils.DirExists(srv.tempDirPath(ctxId))
}
