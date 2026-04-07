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
	log.Println("📦 Application created")

	// ✅ Запуск (должен блокировать)
	log.Println("▶️ Calling application.Run()...")
	if err := application.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}

	log.Println("⚠️ application.Run() returned")
}
