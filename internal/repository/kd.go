package repository

import (
	"context"
	"fmt"
	"vault-exporter/internal/domain"

	"github.com/jackc/pgx/v5"
)

func (repo *ksRepository) AddKdRow(ctx context.Context, data *KdRowDTO, tx pgx.Tx) error {
	var newRel int64

	var idTypeKD int64
	switch data.TypeKD {
	case domain.AssemblyDoc:
		idTypeKD = 10000001316
	case domain.Specification:
		idTypeKD = 10000001315
	case domain.TecDrawing:
		idTypeKD = 10000001321
	case domain.WithoutDoc:
		idTypeKD = 10000001356

	}

	res, err := tx.Query(ctx, `INSERT INTO orl (id_obj_own, id_prmt, id_obj_mem) 
		VALUES ($1, $2, $3) RETURNING id_obj_rlt`, data.Id, 183533, idTypeKD)
	if err != nil {
		return fmt.Errorf("can't add to kd: %w id = %d, file = %s", err, data.Id, *data.FileName)
	}
	defer res.Close()

	res.Next()
	err = res.Scan(&newRel)
	if err != nil {
		return fmt.Errorf("can't add to kd: %w id = %d, file = %s", err, data.Id, *data.FileName)
	}
	res.Close()

	if data.TypeKD == domain.WithoutDoc {
		_, err = tx.Exec(ctx, `INSERT INTO orc (id_obj_rlt, discription, discript2, pln, cnt, id_obj) 
			VALUES ($1, $2, $3, $4, $5, $6)`, newRel, data.Code, data.Name, 1, 1, data.Material)
	} else {
		var cmdText string
		if data.IsPdf {
			cmdText = `INSERT INTO orc (id_obj_rlt, discription, discript2, id_dct, pln, cnt, discript4, id_obj) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
		} else {
			cmdText = `INSERT INTO orc (id_obj_rlt, discription, discript2, id_dct, pln, cnt, discript3, id_obj) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
		}

		_, err = tx.Exec(ctx, cmdText, newRel, data.Code, data.Name, 10000001329, 1, 1, *data.FileName, data.Material)
	}
	if err != nil {
		return fmt.Errorf("can't add to kd: %w id = %d, file = %s", err, data.Id, *data.FileName)
	}

	return nil
}

func (repo *ksRepository) AddMainKd(ctx context.Context, data *KdMainDTO, tx pgx.Tx) error {

	var idPrmt int64

	switch data.TypeFile {
	case domain.KD:
		idPrmt = 1004002
	case domain.MachineFile:
		idPrmt = 13037934
	default:
		return fmt.Errorf("wrong type file when add to main kd")
	}

	_, err := tx.Exec(ctx, `INSERT INTO op (id_obj, id_prmt, vl) 
		VALUES ($1, $2, ($3)::text)`, data.Id, idPrmt, data.FileName)
	if err != nil {
		return fmt.Errorf("can't add to main kd: %w id = %d, file = %d", err, data.Id, data.FileName)
	}

	return nil
}
