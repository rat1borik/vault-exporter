package repository

import (
	"context"
	"fmt"
	"strconv"
	"vault-exporter/internal/domain"
	"vault-exporter/internal/infrastructure"

	"github.com/jackc/pgx/v5"
)

func creationError(place string, code string, wrapped error) error {
	if wrapped != nil {
		return fmt.Errorf("error while creating izd (%s): %w, izd = %s", place, wrapped, code)
	}

	return fmt.Errorf("error while creating izd (%s) = %s", place, code)
}

func (repo *ksRepository) CreateIzd(ctx context.Context, options *IzdCreationOptions, tx pgx.Tx) (int64, error) {
	var newIzdId int64
	res, err := tx.Query(ctx, `INSERT INTO o (act, code, id_grown, id_cls, name, code_cond, id_obj_razrab) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id_obj`, 1, options.Code, IdGrown, IdClsIzd, options.Name, options.CodeName, IdRazrabKD)
	if err != nil {
		return 0, creationError("inserting in o", options.Code, err)
	}
	defer res.Close()

	res.Next()
	err = res.Scan(&newIzdId)
	if err != nil {
		return 0, creationError("reading id", options.Code, err)
	}
	res.Close()

	if err := fillIzdHeaderParameters(ctx, options, tx, newIzdId); err != nil {
		return 0, creationError("filling header", options.Code, err)
	}

	if err := createObjectManagement(ctx, tx, newIzdId, 175723, 1, "Загружено автоматически из Autodesk Vault"); err != nil {
		return 0, creationError("creating object management", options.Code, err)
	}

	return newIzdId, nil
}

func fillIzdHeaderParameters(ctx context.Context, options *IzdCreationOptions, tx pgx.Tx, id int64) error {
	batch := infrastructure.NewCountingBatch()
	// Заполняем группу
	batch.Queue(`INSERT INTO orl (id_obj_own, id_prmt, id_obj_mem) 
		VALUES ($1, $2, $3)`, id, 1000101, options.GroupId)

	// Заполняем раздел спецификации
	batch.Queue(`INSERT INTO orl (id_obj_own, id_prmt, id_obj_mem) 
		VALUES ($1, $2, $3)`, id, 175478, options.SpecDivisionId)

	// Заполняем единицу измерения
	batch.Queue(`INSERT INTO orl (id_obj_own, id_prmt, id_obj_mem) 
		VALUES ($1, $2, $3)`, id, 175486, options.UnitsId)

	if options.Weight != nil {
		weight := strconv.FormatFloat(*options.Weight, 'f', -1, 64)
		// Заполняем массу
		batch.Queue(`INSERT INTO op (id_obj, id_prmt, vl, id_dct_edizm) 
		VALUES ($1, $2, $3, $4)`, id, 175483, weight, domain.Kg)
	}

	// Заполняем вид
	batch.Queue(`INSERT INTO OC (ID_OBJ, ID_CLS, ID_PRMT) 
		VALUES ($1, $2, $3)`, id, IdClsIzd, 82006)

	if err := batch.ExecTx(ctx, tx); err != nil {
		return fmt.Errorf("can't fill parameters: %w", err)
	}

	return nil
}

func createObjectManagement(ctx context.Context, tx pgx.Tx, id int64, shellPrmt int64, stage int, description string) error {
	// shell
	res, err := tx.Query(ctx, "INSERT INTO SHELL (code, id_prmt) VALUES ($1, $2) RETURNING id_shell", "-", shellPrmt)
	if err != nil {
		return fmt.Errorf("can't create object management: %w", err)
	}
	defer res.Close()

	var idShell int64

	res.Next()
	err = res.Scan(&idShell)
	if err != nil {
		return fmt.Errorf("can't create object management: %w", err)
	}
	res.Close()

	// shell_move
	_, err = tx.Exec(ctx, `INSERT INTO SHELL_MOVE (ID_OBJ_STTS, D1, ID_PRSN, ID_SHELL, REASON_RET)
  					          values ($1, current_timestamp, $2, $3, $4)`, stage, IdRazrab, idShell, description)
	if err != nil {
		return fmt.Errorf("can't create object management: %w", err)
	}

	// shell_object
	_, err = tx.Exec(ctx, "INSERT INTO SHELL_OBJECT (id_obj, id_shell) VALUES ($1, $2)", id, idShell)
	if err != nil {
		return fmt.Errorf("can't create object management: %w", err)
	}

	return nil
}

func (repo *ksRepository) AddToAssembly(ctx context.Context, data *AddToAssemblyRepoDTO, tx pgx.Tx) error {
	var newRel int64

	res, err := tx.Query(ctx, `INSERT INTO orl (id_obj_own, id_prmt, id_obj_mem) 
		VALUES ($1, $2, $3) RETURNING id_obj_rlt`, data.ParentId, 82013, data.Id)
	if err != nil {
		return fmt.Errorf("can't add to assembly: %w id = %d, parentId = %d", err, data.Id, data.ParentId)
	}
	defer res.Close()

	res.Next()
	err = res.Scan(&newRel)
	if err != nil {
		return fmt.Errorf("can't add to assembly: %w id = %d, parentId = %d", err, data.Id, data.ParentId)
	}
	res.Close()

	_, err = tx.Exec(ctx, `INSERT INTO orc (id_obj_rlt, id_obj_from, id_dct, id_dct_edizm, cnt, discript5) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id_obj_rlt`, newRel, IdYes, IdYes, domain.Piece, data.Quantity, strconv.Itoa(data.Position))
	if err != nil {
		return fmt.Errorf("can't add to assembly: %w id = %d, parentId = %d", err, data.Id, data.ParentId)
	}

	return nil
}
