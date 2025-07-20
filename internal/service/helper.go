package service

import (
	"strconv"
	"vault-exporter/internal/domain"
)

// Модель добавления нового изделия в сборку
type AddToAssemblyDTO struct {
	ParentId int64
	Id       int64
	Quantity int
	Position int
}

// Модель добавления КД в новое изделие
type KDOptions struct {
	Id           int64
	Positions    []KDPosition
	MaterialName *string
	SpecDivision domain.SpecDivision
	Name         string
	Code         string
}

type KDPosition struct {
	FileName string // Название файла КД
}

type VaultBuilt []*domain.VaultItem

func (a VaultBuilt) Len() int      { return len(a) }
func (a VaultBuilt) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a VaultBuilt) Less(i, j int) bool {
	if a[i].PositionNum == nil || a[j].PositionNum == nil {
		return false
	}

	v1, _ := strconv.Atoi(*a[i].PositionNum)
	v2, _ := strconv.Atoi(*a[j].PositionNum)

	return v1 < v2
}
