package domain

import (
	"strings"
	"vault-exporter/internal/utils"
)

type FileType int

const (
	KD          FileType = iota // Конструкторская документация
	MachineFile                 // Файл для станка
	AnotherFile
)

type KDType int

const (
	AssemblyDoc   KDType = iota // Сборочная документация - СБ
	Specification               // Спецификация - СП
	TecDrawing                  // Чертежная документация - ЧД
	AnotherDoc                  // Другой документ
	WithoutDoc                  // Без документа
)

type PaperFormat int

const (
	A4 PaperFormat = iota
)

// Информация, выявленная из имени файла
type FileProperties struct {
	Type   FileType
	TypeKD KDType
	Format PaperFormat
	Ext    string
	Name   string
}

func ParseFilename(fileName string, sd SpecDivision) (*FileProperties, error) {
	idx := strings.LastIndex(fileName, ".")
	if idx == -1 {
		return nil, utils.UserErrorf("некорректное имя файла - %s", fileName)
	}

	name := fileName[:idx]
	ext := fileName[idx+1:]

	fType := fileType(ext)

	if fType == MachineFile {
		return &FileProperties{
			Type:   fType,
			TypeKD: AnotherDoc,
			Format: A4,
			Name:   name,
			Ext:    ext,
		}, nil
	}

	var kdType KDType

	switch sd {
	case Assembly:
		if strings.Contains(name, " СП ") {
			kdType = Specification
		} else {
			kdType = AssemblyDoc
		}
	case Part:
		kdType = TecDrawing
	default:
		kdType = AnotherDoc
	}

	return &FileProperties{
		Type:   fType,
		TypeKD: kdType,
		Format: A4,
		Name:   name,
		Ext:    ext,
	}, nil
}

func fileType(ext string) FileType {
	ext2Type := map[string]FileType{
		"pdf": KD,
		"dxf": MachineFile,
	}

	if val, ok := ext2Type[ext]; ok {
		return val
	}

	return AnotherFile
}
