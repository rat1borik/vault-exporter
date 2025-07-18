package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type KSRepository interface {
	CreateIzd(ctx context.Context, options *IzdCreationOptions, tx pgx.Tx) (int64, error)
	AddToAssembly(ctx context.Context, data *AddToAssemblyRepoDTO, tx pgx.Tx) error
	FindNm(context.Context, *FindNmDTO, pgx.Tx) (int64, error)
}

type ksRepository struct {
	Db *pgxpool.Pool
}

func NewKSRepository(db *pgxpool.Pool) KSRepository {
	return &ksRepository{Db: db}
}
