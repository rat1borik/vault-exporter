// Package service содержит слой сервисов
package service

import (
	"vault-exporter/internal/config"
	"vault-exporter/internal/model"
	"vault-exporter/internal/repository"
)

type IzdCreatorService interface {
	CreateIzd(item *model.VaultItem) (int64, error)
}

type izdCreatorService struct {
	cfg  *config.ServerConfig
	repo *repository.KSRepository
}

func NewIzdCreatorService(cfg *config.ServerConfig, repo *repository.KSRepository) IzdCreatorService {
	return &izdCreatorService{cfg: cfg, repo: repo}
}

func (svc *izdCreatorService) CreateIzd(item *model.VaultItem) (int64, error) {

	return 0, nil
}
