package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"vault-exporter/internal/config"
	"vault-exporter/internal/domain"
	"vault-exporter/internal/repository"
	"vault-exporter/internal/utils"

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

	tx, err := svc.db.Begin(ctx)
	if err != nil {
		return []error{fmt.Errorf("can't start transaction: %w", err)}
	}
	defer tx.Rollback(ctx)

	var root *domain.VaultItem

	for i, el := range val {
		if el.ParentId == nil {
			root = &val[i]
			break
		}
	}

	errs := utils.NewUserErrorCollection()

	created := make(map[int64]struct{}, len(val))

	var createAndFill func(node *domain.VaultItem) *int64
	createAndFill = func(node *domain.VaultItem) *int64 {
		// Проверяем на ацикличность
		if _, ok := created[node.Id]; ok {
			errs.Add(utils.UserErrorf("обнаружено зацикливание в узле %s %s", node.Title, node.PartNumber), !svc.cfg.IsProduction)
			return nil
		}

		created[node.Id] = struct{}{}

		idParent, err := svc.izdCreatorSvc.CreateOrFindNm(ctx, node, tx)
		if err != nil {
			errs.Add(err, !svc.cfg.IsProduction)
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
					errs.Add(utils.UserErrorf("не удается определить позицию в сборке %s %s", node.Title, node.PartNumber), !svc.cfg.IsProduction)
					continue

				}
				err = svc.izdCreatorSvc.AddToAssembly(ctx, &AddToAssemblyDTO{
					ParentId: idParent,
					Id:       *id,
					Quantity: int(*val[i].Quant),
					Position: pos,
				}, tx)
				if err != nil {
					var dummy *utils.UserError
					if errors.As(err, &dummy) || !svc.cfg.IsProduction {
						errs.Add(err, !svc.cfg.IsProduction)
					}

					continue
				}
			}
		}

		return &idParent
	}

	createAndFill(root)

	errsFinal := errs.Collection()
	if errsFinal == nil {
		if err := tx.Commit(ctx); err != nil {
			svc.fileGetterSvc.ClearTempFolder(ctx.Value(utils.CtxProcId).(string))
			return []error{err}
		}
		svc.fileGetterSvc.CommitFiles(ctx.Value(utils.CtxProcId).(string))
	}
	svc.fileGetterSvc.ClearTempFolder(ctx.Value(utils.CtxProcId).(string))
	return errsFinal
}
