package service

import (
	"context"
	"fmt"
	"strconv"
	"vault-exporter/internal/config"
	"vault-exporter/internal/domain"
	"vault-exporter/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Оркестрирует полный процесс импорта
type ImportProcessorService interface {
	Import(context.Context, []domain.VaultItem) []error
}

type importProcessorService struct {
	cfg           *config.ServerConfig
	db            *pgxpool.Pool
	repo          repository.KSRepository
	fileGetterSvc FileGetterService
	izdCreatorSvc IzdCreatorService
}

func NewImportProcessorService(cfg *config.ServerConfig, db *pgxpool.Pool, fileGetterSvc FileGetterService, izdCreatorService IzdCreatorService, repo repository.KSRepository) ImportProcessorService {
	return &importProcessorService{
		cfg:           cfg,
		db:            db,
		repo:          repo,
		fileGetterSvc: fileGetterSvc,
		izdCreatorSvc: izdCreatorService}
}

func (svc importProcessorService) Import(ctx context.Context, val []domain.VaultItem) []error {
	if len(val) == 0 {
		return nil
	}

	tx, err := svc.db.Begin(context.Background())
	if err != nil {
		return []error{fmt.Errorf("can't start transaction: %w", err)}
	}
	defer tx.Rollback(context.Background())

	var root *domain.VaultItem

	for i, el := range val {
		if el.ParentId == nil {
			root = &val[i]
			break
		}
	}

	errs := make([]error, 0)

	created := make(map[int64]struct{}, len(val))

	var createAndFill func(node *domain.VaultItem) *int64
	createAndFill = func(node *domain.VaultItem) *int64 {
		// Проверяем на ацикличность
		if _, ok := created[node.Id]; ok {
			errs = append(errs, fmt.Errorf("cycle detected at nodeId = %d", node.Id))
			return nil
		}

		created[node.Id] = struct{}{}

		idParent, err := svc.izdCreatorSvc.CreateIzd(ctx, node, tx)
		if err != nil {
			errs = append(errs, err)
			return nil
		}

		// Ищем детей
		for i := range val {
			if val[i].ParentId != nil && *val[i].ParentId == node.Id {
				id := createAndFill(&val[i])
				if id == nil {
					continue
				}

				// Тут добавляем в состав
				pos, err := strconv.Atoi(*val[i].PositionNum)
				if err != nil {
					errs = append(errs, fmt.Errorf("can't define position in assembly, nodeId = %d", val[i].Id))
					continue

				}
				err = svc.izdCreatorSvc.AddToAssembly(ctx, &AddToAssemblyDTO{
					ParentId: idParent,
					Id:       *id,
					Quantity: int(*val[i].Quant),
					Position: pos,
				}, tx)
				if err != nil {
					errs = append(errs, err)
					continue
				}
			}
		}

		return &idParent
	}

	createAndFill(root)

	if len(errs) != 0 {
		return errs
	}

	tx.Commit(context.Background())
	return nil
}
