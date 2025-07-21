BINARY_NAME = ve.exe

# Путь к main пакету
MAIN_PACKAGE = ./cmd/server

# Папка для билдов
BUILD_DIR = build

# Опции компиляции (опционально)
GOFLAGS = -ldflags="-s -w -X 'main.AppEnv=production'"

.PHONY: all build clean run fmt lint tidy

# Компилируем бинарь
build:
	go build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)

# Удаляем build
clean:
	del /Q $(BUILD_DIR)\* || echo "Nothing to clean."

# Запускаем приложение (для разработки)
run:
	go run $(MAIN_PACKAGE)

# Запускаем приложение (для разработки)
run-race:
	go run -race $(MAIN_PACKAGE)

# Форматируем код
fmt:
	go fmt ./...

# Проверяем линтинг
lint:
	golangci-lint run

# Устанавливаем зависимости
tidy:
	go mod tidy

# Сборка по умолчанию
all: tidy fmt build