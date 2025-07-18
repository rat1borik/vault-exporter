package domain

import (
	"errors"
)

// Разделы спецификаций
type SpecDivision int64

const (
	Assembly SpecDivision = 10000001424 // Сборка
	Part     SpecDivision = 10000001425 // Деталь
	Another  SpecDivision = 0
)

func DefSpecDivision(val string) SpecDivision {
	specDivisionMap := map[string]SpecDivision{
		"Part":     Part,
		"Assembly": Assembly,
	}

	if res, ok := specDivisionMap[val]; ok {
		return res
	}

	return Another
}

// Единицы измерения
type Unit int64

const (
	Piece Unit = 10000001206 // Штука
	Kg    Unit = 10000001225 // Килограмм
)

func DefUnit(val int) (Unit, error) {
	unitsMap := map[int]Unit{
		1: Piece,
	}

	if res, ok := unitsMap[val]; ok {
		return res, nil
	}

	return 0, errors.New("can't define unit")
}

// Группы
type IzdGroup int64

const (
	MK IzdGroup = 130000104985 // Металлоконструкции
)
