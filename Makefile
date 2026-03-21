.PHONY: help build build-dev build-release package run test clean resources generate

# Переменные
APP_NAME := TrayApp
MODULE := fyne-tray-app
CMD_PATH := ./cmd/app
ICON_PATH := resources/icons/icon-on.ico
OUTPUT_DIR := build

# Цвета для вывода (работают в PowerShell/bash)
COLOR_RESET := \033[0m
COLOR_GREEN := \033[32m
COLOR_BLUE := \033[34m
COLOR_YELLOW := \033[33m

help:
	@echo "$(COLOR_BLUE)Доступные команды:$(COLOR_RESET)"
	@echo "  make build-dev      - Сборка для разработки (с консолью)"
	@echo "  make build-release  - Сборка релиза (без консоли, windowsgui)"
	@echo "  make package        - Упаковка через fyne package"
	@echo "  make run            - Запуск без сборки"
	@echo "  make test           - Запуск тестов"
	@echo "  make clean          - Очистка кэша и билдов"
	@echo "  make resources      - Генерация ресурсов (fyne bundle)"
	@echo "  make generate       - Генерация всех автогенерируемых файлов"

# Генерация ресурсов (иконки и др.)
resources:
	@echo "$(COLOR_GREEN)→ Генерация ресурсов...$(COLOR_RESET)"
	@if exist resources\icons ( \
		fyne bundle -o internal/resources/icons/icons.go -package icons resources/icons/ \
	) else ( \
		echo "⚠️ Папка resources/icons не найдена, пропускаем" \
	)

# Генерация всех ресурсов (go generate)
generate: resources
	@echo "$(COLOR_GREEN)→ Запуск go generate...$(COLOR_RESET)"
	go generate ./...

# Сборка для разработки (с консолью, без оптимизаций)
build-dev: generate
	@echo "$(COLOR_GREEN)→ Сборка для разработки...$(COLOR_RESET)"
	@mkdir -p $(OUTPUT_DIR)
	set CGO_ENABLED=1 && go build -v -o $(OUTPUT_DIR)/$(APP_NAME)-dev.exe $(CMD_PATH)
	@echo "$(COLOR_GREEN)✓ Готово: $(OUTPUT_DIR)/$(APP_NAME)-dev.exe$(COLOR_RESET)"

# Сборка релиза (без консоли, с оптимизациями)
build-release: generate
	@echo "$(COLOR_GREEN)→ Сборка релиза...$(COLOR_RESET)"
	@mkdir -p $(OUTPUT_DIR)
	set CGO_ENABLED=1 && go build -ldflags="-H=windowsgui -s -w" -o $(OUTPUT_DIR)/$(APP_NAME).exe $(CMD_PATH)
	@echo "$(COLOR_GREEN)✓ Готово: $(OUTPUT_DIR)/$(APP_NAME).exe$(COLOR_RESET)"

# Упаковка через fyne package (с иконкой, метаданными)
package: generate
	@echo "$(COLOR_GREEN)→ Упаковка приложения...$(COLOR_RESET)"
	@mkdir -p $(OUTPUT_DIR)
	fyne package -os windows -icon $(ICON_PATH) -appID com.example.trayapp -name $(APP_NAME) -release
	@echo "$(COLOR_GREEN)✓ Готово: $(OUTPUT_DIR)/$(APP_NAME).exe$(COLOR_RESET)"

# Запуск без сборки (для разработки)
run: generate
	@echo "$(COLOR_GREEN)→ Запуск приложения...$(COLOR_RESET)"
	set CGO_ENABLED=1 && go run $(CMD_PATH)

# Запуск тестов
test:
	@echo "$(COLOR_GREEN)→ Запуск тестов...$(COLOR_RESET)"
	go test -v -cover ./...

# Запуск тестов с покрытием
test-coverage:
	@echo "$(COLOR_GREEN)→ Запуск тестов с покрытием...$(COLOR_RESET)"
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(COLOR_GREEN)✓ Отчёт: coverage.html$(COLOR_RESET)"

# Очистка
clean:
	@echo "$(COLOR_YELLOW)→ Очистка...$(COLOR_RESET)"
	go clean -cache -modcache
	rm -rf $(OUTPUT_DIR)/
	rm -f coverage.out coverage.html
	@echo "$(COLOR_GREEN)✓ Очистка завершена$(COLOR_RESET)"

# Установка зависимостей
deps:
	@echo "$(COLOR_GREEN)→ Установка зависимостей...$(COLOR_RESET)"
	go mod download
	go mod tidy
	@echo "$(COLOR_GREEN)✓ Зависимости установлены$(COLOR_RESET)"

# Установка утилит (fyne, goimports и др.)
tools:
	@echo "$(COLOR_GREEN)→ Установка инструментов...$(COLOR_RESET)"
	go install fyne.io/fyne/v2/cmd/fyne@latest
	go install golang.org/x/tools/cmd/goimports@latest
	@echo "$(COLOR_GREEN)✓ Инструменты установлены$(COLOR_RESET)"

# Полная пересборка
rebuild: clean build-release
	@echo "$(COLOR_GREEN)✓ Пересборка завершена$(COLOR_RESET)"

# Сборка для всех платформ (кросс-компиляция)
build-all: generate
	@echo "$(COLOR_GREEN)→ Сборка для всех платформ...$(COLOR_RESET)"
	@mkdir -p $(OUTPUT_DIR)/windows
	@mkdir -p $(OUTPUT_DIR)/linux
	@mkdir -p $(OUTPUT_DIR)/darwin

	# Windows
	set CGO_ENABLED=1 && GOOS=windows GOARCH=amd64 go build -ldflags="-H=windowsgui -s -w" -o $(OUTPUT_DIR)/windows/$(APP_NAME).exe $(CMD_PATH)

	# Linux
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(OUTPUT_DIR)/linux/$(APP_NAME) $(CMD_PATH)

	# macOS
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $(OUTPUT_DIR)/darwin/$(APP_NAME) $(CMD_PATH)

	@echo "$(COLOR_GREEN)✓ Сборка завершена для всех платформ$(COLOR_RESET)"