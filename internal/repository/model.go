package repository

import "vault-exporter/internal/domain"

// Голые данные для загона в БД
type IzdCreationOptions struct {
	Code            string              // Обозначение
	Name            string              // Наименование
	CodeName        string              // Условное обозначение
	SpecDivisionId  domain.SpecDivision // Раздел спецификации
	UnitsId         domain.Unit         // Единица измерения
	GroupId         domain.IzdGroup     // Группа
	Weight          float64             // Масса
	FileMachineName string              // Файл для станка
	MainFileName    string              // Файл основной
}

const IdRazrab int64 = 13004684
const IdRazrabKD int64 = 130000583353 // Разраб для всего

const IdGrown int64 = 768354 // Организация под все нужды
const IdClsIzd int32 = 10    // Класс ИЗД
