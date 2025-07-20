// Package service содержит слой сервисов
package service

import (
	"context"
	"fmt"
	"sync"
	"vault-exporter/internal/config"
	"vault-exporter/internal/domain"
	"vault-exporter/internal/repository"
	"vault-exporter/internal/utils"

	"github.com/jackc/pgx/v5"
)

type IzdCreatorService interface {
	CreateOrFindNm(ctx context.Context, item *domain.VaultItem, tx pgx.Tx) (int64, error)
	AddToAssembly(ctx context.Context, data *AddToAssemblyDTO, tx pgx.Tx) error
	AddKD(ctx context.Context, data *KDOptions, tx pgx.Tx) error
}

type izdCreatorService struct {
	cfg           *config.ServerConfig
	repo          repository.KSRepository
	fileGetterSvc FileGetterService
}

func NewIzdCreatorService(cfg *config.ServerConfig, fileGetterSvc FileGetterService, repo repository.KSRepository) IzdCreatorService {
	return &izdCreatorService{cfg: cfg,
		repo:          repo,
		fileGetterSvc: fileGetterSvc}
}

func (svc *izdCreatorService) CreateOrFindNm(ctx context.Context, item *domain.VaultItem, tx pgx.Tx) (int64, error) {
	spec := domain.DefSpecDivision(item.CatSystemName)

	// Не производится - нужно найти
	if spec == domain.AnotherSpec {
		idNm, err := svc.repo.FindNm(ctx, &repository.FindNmDTO{
			Name: item.PartNumber,
		}, tx)

		if err != nil {
			return 0, err
		}

		return idNm, nil
	}

	var fileError error

	wgFiles := sync.WaitGroup{}
	files := utils.NewSafeSlice[KDPosition]()

	for i := range item.Files {
		wgFiles.Add(1)
		go func() {
			defer wgFiles.Done()
			newFile, err := svc.fileGetterSvc.LoadFile(&item.Files[i], ctx.Value("proc_id").(string))
			if err != nil && fileError == nil {
				fileError = err
			}

			files.Append(KDPosition{
				FileName: newFile,
			})
		}()
	}

	unit, err := domain.DefUnit(*item.UnitID)
	if err != nil {
		return 0, fmt.Errorf("can't define unit izd: %w", err)
	}

	props := propertiesMap(item.Properties)

	var weight *float64

	if val, ok := props[122]; ok {
		if val != nil {
			w := val.(float64)
			weight = &w
		}
	}

	opts := &repository.IzdCreationOptions{
		Code:           item.PartNumber,
		Name:           item.Title,
		CodeName:       item.Title,
		SpecDivisionId: spec,
		UnitsId:        unit,
		GroupId:        domain.MK,
		Weight:         weight,
	}

	id, err := svc.repo.CreateIzd(ctx, opts, tx)
	if err != nil {
		return 0, err
	}

	var material *string
	if val, ok := props[100]; ok {
		s := val.(string)
		material = &s
	}

	wgFiles.Wait()
	if fileError != nil {
		return 0, fileError
	}

	if err := svc.AddKD(ctx, &KDOptions{
		Id:           id,
		SpecDivision: spec,
		Name:         item.Title,
		Code:         item.PartNumber,
		Positions:    files.Items(),
		MaterialName: material,
	}, tx); err != nil {
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

func (svc *izdCreatorService) AddKD(ctx context.Context, data *KDOptions, tx pgx.Tx) error {
	var idMat *int64

	writtenKdMain := false
	writtenMachineMain := false

	if data.MaterialName != nil && data.SpecDivision == domain.Part {
		id, err := svc.repo.FindNm(ctx, &repository.FindNmDTO{Name: *data.MaterialName}, tx)
		if err != nil {
			return fmt.Errorf("ошибка добавления КД (%s): %w", data.Code, err)
		}

		idMat = &id
	}

	if len(data.Positions) == 0 {
		if err := svc.repo.AddKdRow(ctx, &repository.KdRowDTO{
			FileName: nil,
			TypeKD:   domain.WithoutDoc,
			Material: idMat,
			Id:       data.Id,
			Format:   domain.A4,
			Name:     data.Name,
			Code:     data.Code,
		}, tx); err != nil {
			return fmt.Errorf("ошибка добавления КД (%s): %w", data.Code, err)
		}

		return nil
	}

	for _, pos := range data.Positions {
		parsed, err := domain.ParseFilename(pos.FileName, data.SpecDivision)
		if err != nil {
			return fmt.Errorf("ошибка определения информации в файле %s (%s): %w", pos.FileName, data.Code, err)
		}

		if !writtenKdMain && parsed.Type == domain.KD {
			if err := svc.repo.AddMainKd(ctx, &repository.KdMainDTO{
				Id:       data.Id,
				TypeFile: parsed.Type,
				FileName: &pos.FileName,
			}, tx); err != nil {
				return fmt.Errorf("can't write main kd %s (%s): %w", pos.FileName, data.Code, err)
			}

			writtenKdMain = true
		}

		if !writtenMachineMain && parsed.Type == domain.MachineFile {
			if err := svc.repo.AddMainKd(ctx, &repository.KdMainDTO{
				Id:       data.Id,
				TypeFile: parsed.Type,
				FileName: &pos.FileName,
			}, tx); err != nil {
				return fmt.Errorf("can't write main kd %s (%s): %w", pos.FileName, data.Code, err)
			}

			writtenMachineMain = true
		}

		if parsed.Type == domain.MachineFile {
			// Файл для станка не заносится в таблицу
			continue
		}

		if err := svc.repo.AddKdRow(ctx, &repository.KdRowDTO{
			FileName: &pos.FileName,
			IsPdf:    parsed.Ext == "pdf",
			TypeKD:   parsed.TypeKD,
			Material: idMat,
			Id:       data.Id,
			Format:   domain.A4,
			Name:     data.Name,
			Code:     data.Code,
		}, tx); err != nil {
			return fmt.Errorf("ошибка добавления КД (%s): %w", data.Code, err)
		}
	}

	return nil
}

func filesToKDPositions(val []domain.VaultFile) []KDPosition {
	res := make([]KDPosition, 0, len(val))

	for i := range val {
		res = append(res, KDPosition{FileName: val[i].FileName})
	}

	return res
}
