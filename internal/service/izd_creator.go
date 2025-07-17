// Package service содержит слой сервисов
package service

import (
	"context"
	"fmt"
	"vault-exporter/internal/config"
	"vault-exporter/internal/domain"
	"vault-exporter/internal/repository"

	"github.com/jackc/pgx/v5"
)

type IzdCreatorService interface {
	CreateIzd(ctx context.Context, item *domain.VaultItem, tx pgx.Tx) (int64, error)
	AddToAssembly(ctx context.Context, data *AddToAssemblyDTO, tx pgx.Tx) error
}

type izdCreatorService struct {
	cfg  *config.ServerConfig
	repo repository.KSRepository
}

func NewIzdCreatorService(cfg *config.ServerConfig, repo repository.KSRepository) IzdCreatorService {
	return &izdCreatorService{cfg: cfg, repo: repo}
}

func (svc *izdCreatorService) CreateIzd(ctx context.Context, item *domain.VaultItem, tx pgx.Tx) (int64, error) {
	spec, err := domain.DefSpecDivision(item.CatSystemName)
	if err != nil {
		return 0, fmt.Errorf("can't define spec izd: %w", err)
	}

	unit, err := domain.DefUnit(*item.UnitID)
	if err != nil {
		return 0, fmt.Errorf("can't define unit izd: %w", err)
	}

	props := propertiesMap(item.Properties)

	opts := &repository.IzdCreationOptions{
		Code:           item.PartNumber,
		Name:           item.Title,
		CodeName:       item.Title,
		SpecDivisionId: spec,
		UnitsId:        unit,
		GroupId:        domain.MK,
		Weight:         props[122].(float64),
	}

	id, err := svc.repo.CreateIzd(ctx, opts, tx)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func propertiesMap(val []domain.VaultProperty) map[int]interface{} {
	res := make(map[int]interface{}, len(val))

	for _, v := range val {
		res[v.PropDefId] = v.Val
	}

	return res
}

func (svc *izdCreatorService) AddToAssembly(ctx context.Context, data *AddToAssemblyDTO, tx pgx.Tx) error {
	// TODO: какие-то валидации и проч. вещи
	return svc.repo.AddToAssembly(ctx, &repository.AddToAssemblyRepoDTO{
		ParentId: data.ParentId,
		Id:       data.Id,
		Quantity: data.Quantity,
		Position: data.Position,
	}, tx)

}
