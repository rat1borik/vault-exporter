package service

import (
	"database/sql"
	"fmt"
	"vault-exporter/internal/config"
	"vault-exporter/internal/domain"
	"vault-exporter/internal/repository"
)

// Оркестрирует полный процесс импорта
type ImportProcessorService interface {
	Import([]domain.VaultItem) error
}

type importProcessorService struct {
	cfg           *config.ServerConfig
	db            *sql.DB
	repo          repository.KSRepository
	fileGetterSvc FileGetterService
	izdCreatorSvc IzdCreatorService
}

func NewImportProcessorService(cfg *config.ServerConfig, db *sql.DB, fileGetterSvc FileGetterService, izdCreatorService IzdCreatorService, repo repository.KSRepository) ImportProcessorService {
	return &importProcessorService{
		cfg:           cfg,
		db:            db,
		repo:          repo,
		fileGetterSvc: fileGetterSvc,
		izdCreatorSvc: izdCreatorService}
}

func (svc importProcessorService) Import(val []domain.VaultItem) error {
	if len(val) == 0 {
		return nil
	}

	tx, err := svc.db.Begin()
	if err != nil {
		return fmt.Errorf("can't start transaction")
	}
	defer tx.Rollback()

	svc.izdCreatorSvc.CreateIzd(&val[0], tx)

	tx.Commit()
	return nil
}
