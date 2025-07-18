package repository

import (
	"context"
	"vault-exporter/internal/utils"

	"github.com/jackc/pgx/v5"
)

func (repo *ksRepository) FindNm(ctx context.Context, options *FindNmDTO, tx pgx.Tx) (int64, error) {
	var idNm int64
	res := tx.QueryRow(ctx, "SELECT id_obj FROM o WHERE id_cls = $1 AND act != 2 AND name = $2", 66533, options.Name)
	if err := res.Scan(&idNm); err == nil {
		return idNm, nil
	} else if err != pgx.ErrNoRows {
		return 0, err
	}

	// Не нашли - идем в синонимы

	// TODO: идти в синонимы

	return 0, utils.UserErrorf(`не удается найти подходящую номенклатуру для "%s"`, options.Name)
}
