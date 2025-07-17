// Package utils - полезные функции обертки
package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err != nil
}

func DirExists(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}

	return stat.IsDir()
}

func EnsureDir(path string) error {
	if DirExists(path) {
		return nil
	}

	err := os.MkdirAll(path, 0755)

	if err != nil {
		return fmt.Errorf("can't ensure dir: %w", err)
	}

	return nil
}

func ClearDir(path string, self bool) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("can't read folder: %w", err)
	}

	for _, entry := range entries {
		entryPath := filepath.Join(path, entry.Name())
		err := os.RemoveAll(entryPath) // удаляет как файлы, так и вложенные папки
		if err != nil {
			return fmt.Errorf("can't remove %s: %w", entryPath, err)
		}
	}

	if self {
		err := os.Remove(path)
		if err != nil {
			return fmt.Errorf("can't remove %s: %w", path, err)
		}
	}

	return nil
}

func CopyFile(srcPath, dstPath string, removeOriginal bool) error {
	// Открыть исходный файл
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("can't openc src: %w", err)
	}
	defer srcFile.Close()

	// Создать целевой файл
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("can't create target: %w", err)
	}
	defer dstFile.Close()

	// Копирование содержимого
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("copy error: %w", err)
	}

	if removeOriginal {
		// Удалить оригинал
		if err := os.Remove(srcPath); err != nil {
			return fmt.Errorf("can't remove original: %w", err)
		}
	}

	return nil
}
