// Package service содержит слой сервисов
package service

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"vault-exporter/internal/config"
	"vault-exporter/internal/domain"
	"vault-exporter/internal/utils"

	"github.com/google/uuid"
)

type FileGetterService interface {
	LoadFile(*domain.VaultFile, string) (string, error)
	ClearTempFolder(string) error
	CommitFiles(string) error
	EnsureUnique(string) bool
}

type fileGetterService struct {
	cfg        *config.ServerConfig
	mu         sync.Mutex
	currentCtx string
}

func NewFileGetterService(cfg *config.ServerConfig) FileGetterService {
	return &fileGetterService{cfg: cfg}
}

// Осуществляет загрузку файла из Vault в специальную директорию КС
func (srv *fileGetterService) LoadFile(file *domain.VaultFile, ctxId string) (string, error) {

	path := fmt.Sprintf("http://%s:%d/api/files?id=%d", srv.cfg.Vault.Host, srv.cfg.Vault.Port, file.Id)

	resp, err := http.Get(path)
	if err != nil {
		return "", fmt.Errorf("сan't perform request to vault: %w, id = %d", err, file.Id)
	}
	defer resp.Body.Close()

	// Пишем во временное место (потом заберем)
	utils.EnsureDir(srv.tempDirPath(ctxId))

	tempName := uuid.New().String()
	dest, err := os.Create(filepath.Join(srv.tempDirPath(ctxId), tempName))
	if err != nil {
		return "", fmt.Errorf("сan't create file: %w, id = %d", err, file.Id)
	}

	// Хеш и запись в файл одновременно
	hasher := sha256.New()
	multiWriter := io.MultiWriter(dest, hasher)

	_, err = io.Copy(multiWriter, resp.Body)
	if err != nil {
		return "", fmt.Errorf("сan't write file: %w, id = %d", err, file.Id)
	}
	dest.Close()

	checksum := hasher.Sum(nil)

	// Захватываем доступ к КС-ной папке для всего контекста
	if srv.currentCtx != ctxId {
		srv.mu.Lock()
		srv.currentCtx = ctxId
	}

	filename := file.FileName

	if utils.FileExists(filepath.Join(srv.cfg.KSFilesPath, file.FileName)) {
		sameFile, err := sameFileExists(filepath.Join(srv.cfg.KSFilesPath, file.FileName), checksum)
		if err != nil {
			return "", fmt.Errorf("сan't ensure same file: %w, id = %d", err, file.Id)
		}

		if !sameFile {
			idx := strings.LastIndex(filename, ".")
			name := filename[:idx]
			ext := filename[idx+1:]
			filename = fmt.Sprintf("%s %s.%s", name, time.Now().Format("02.01.2006 15 04 05"), ext)
		}
	}

	if err = os.Rename(filepath.Join(srv.tempDirPath(ctxId), tempName), filepath.Join(srv.tempDirPath(ctxId), filename)); err != nil {
		return "", err
	}

	return filename, nil
}

func sameFileExists(orig string, newFileHash []byte) (bool, error) {
	if !utils.FileExists(orig) {
		return false, nil
	}

	r, err := os.Open(orig)
	if err != nil {
		return false, err
	}

	hasher := sha256.New()
	io.Copy(hasher, r)

	currentFileHash := hasher.Sum(nil)

	return bytes.Equal(currentFileHash, newFileHash), nil
}

func (srv *fileGetterService) tempDirPath(ctxId string) string {
	return filepath.Join(srv.cfg.TempPath, ctxId)
}

// Чистит временнную папку в случае неудачи / после всех действий
func (srv *fileGetterService) ClearTempFolder(ctxId string) error {
	if ctxId == srv.currentCtx {
		srv.currentCtx = ""
		defer srv.mu.Unlock()
	}
	err := utils.ClearDir(srv.tempDirPath(ctxId), true)
	if err != nil {
		return fmt.Errorf("can't clear temp dir: %w", err)
	}

	return nil
}

// Заливает файлы в финальную папку КС
func (srv *fileGetterService) CommitFiles(ctxId string) error {
	// Захватываем доступ к КС-ной папке для всего контекста
	if srv.currentCtx != ctxId {
		return fmt.Errorf("can't commit not current context")
	}

	defer func() {
		srv.currentCtx = ""
		srv.mu.Unlock()
	}()

	entries, err := os.ReadDir(srv.tempDirPath(ctxId))
	if err != nil {
		return fmt.Errorf("can't commit files: %w", err)
	}

	for _, e := range entries {
		utils.CopyFile(filepath.Join(srv.tempDirPath(ctxId), e.Name()), filepath.Join(srv.cfg.KSFilesPath, e.Name()), false)
	}

	return nil
}

// Проверяет нет ли коллизий по данному ctxId (должен вызываться в самом начале )
func (srv *fileGetterService) EnsureUnique(ctxId string) bool {
	return !utils.DirExists(srv.tempDirPath(ctxId))
}
