package repository

import (
	"database/sql"
)

type KSRepository interface {
	CreateIzd() error
}

type ksRepository struct {
	Db *sql.DB
}

func NewKSRepository(db *sql.DB) KSRepository {
	return &ksRepository{Db: db}
}

func (repo *ksRepository) CreateIzd() error {

	return nil
}
