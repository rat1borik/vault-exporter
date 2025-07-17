package repository

import (
	"database/sql"
	"fmt"
	"vault-exporter/internal/domain"
)

type KSRepository interface {
	CreateIzd(options *IzdCreationOptions, tx *sql.Tx) (int64, error)
}

type ksRepository struct {
	Db *sql.DB
}

func NewKSRepository(db *sql.DB) KSRepository {
	return &ksRepository{Db: db}
}

func creationError(code string, wrapped error) error {
	if wrapped != nil {
		return fmt.Errorf("error while creating izd: %w, izd = %s", wrapped, code)
	}

	return fmt.Errorf("error while creating izd = %s", code)
}

func (repo *ksRepository) CreateIzd(options *IzdCreationOptions, tx *sql.Tx) (int64, error) {
	var newIzdId int64

	res, err := tx.Query(`INSERT INTO o (act, code, id_grown, id_cls, name, code_cond, id_obj_razrab) 
		VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id_obj`, 1, options.Code, IdGrown, IdClsIzd, options.Name, options.CodeName, IdRazrabKD)
	if err != nil {
		return 0, creationError(options.Code, err)
	}
	defer res.Close()

	res.Next()
	err = res.Scan(&newIzdId)
	if err != nil {
		return 0, creationError(options.Code, err)
	}

	if err := fillIzdHeaderParameters(options, tx, newIzdId); err != nil {
		return 0, creationError(options.Code, err)
	}

	if err := createObjectManagement(tx, newIzdId, 175723, 1, "Загружено автоматически из Autodesk Vault"); err != nil {
		return 0, creationError(options.Code, err)
	}

	return newIzdId, nil
}

func fillIzdHeaderParameters(options *IzdCreationOptions, tx *sql.Tx, id int64) error {
	// Заполняем группу
	_, err := tx.Exec(`INSERT INTO orl (id_obj_own, id_prmt, id_obj_mem) 
		VALUES (?, ?, ?)`, id, 1000101, options.GroupId)
	if err != nil {
		return fmt.Errorf("can't fill parameters: %w", err)
	}

	// Заполняем раздел спецификации
	_, err = tx.Exec(`INSERT INTO orl (id_obj_own, id_prmt, id_obj_mem) 
		VALUES (?, ?, ?)`, id, 175478, options.SpecDivisionId)
	if err != nil {
		return fmt.Errorf("can't fill parameters: %w", err)
	}

	// Заполняем единицу измерения
	_, err = tx.Exec(`INSERT INTO orl (id_obj_own, id_prmt, id_obj_mem) 
		VALUES (?, ?, ?)`, id, 175486, options.UnitsId)
	if err != nil {
		return fmt.Errorf("can't fill parameters: %w", err)
	}

	// Заполняем массу
	_, err = tx.Exec(`INSERT INTO op (id_obj_own, id_prmt, vl, id_dct_edizm) 
		VALUES (?, ?, ?, ?)`, id, 175483, options.Weight, domain.Kg)
	if err != nil {
		return fmt.Errorf("can't fill parameters: %w", err)
	}

	// Заполняем вид
	_, err = tx.Exec(`INSERT INTO OC (ID_OBJ, ID_CLS, ID_PRMT) 
		VALUES (?, ?, ?)`, id, IdClsIzd, 82006)
	if err != nil {
		return fmt.Errorf("can't fill parameters: %w", err)
	}

	return nil
}

func createObjectManagement(tx *sql.Tx, id int64, shellPrmt int64, stage int, description string) error {
	// shell
	res, err := tx.Query("INSERT INTO SHELL (code, id_prmt) VALUES (?, ?) RETURNING id_shell", "-", shellPrmt)
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

	// shell_move
	_, err = tx.Exec(`INSERT INTO SHELL_MOVE (ID_OBJ_STTS, D1, ID_PRSN, ID_SHELL, REASON_RET)
  					          values (?, current_timestamp, ?, ?, ?)`, stage, IdRazrab, description)
	if err != nil {
		return fmt.Errorf("can't create object management: %w", err)
	}

	// shell_object
	_, err = tx.Exec("INSERT INTO SHELL_OBJECT (id_obj, id_shell) VALUES (?, ?)", id, idShell)
	if err != nil {
		return fmt.Errorf("can't create object management: %w", err)
	}

	return nil
}
