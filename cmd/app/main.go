package main

import (
	"fyne-tray-app/internal/app"
	"fyne-tray-app/internal/utils"
	"log"
)

func main() {
	log.Println("🚀 Starting Tray App...")

	// Инициализация логгера
	if err := utils.InitLogger(); err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}

	// Создание приложения
	application := app.New()
	log.Println("📦 Application created") // ← Добавьте эту строку!

	// ✅ Запуск (должен блокировать)
	log.Println("▶️ Calling application.Run()...") // ← Добавьте эту строку!
	if err := application.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}

	// Сюда код не дойдёт (Run блокирует)
	log.Println("⚠️ application.Run() returned")
}
