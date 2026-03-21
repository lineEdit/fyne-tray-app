package main

import (
	"fyne-tray-app/internal/app"
	"fyne-tray-app/internal/utils"
	"log"
)

func main() {
	// Инициализация логгера
	if err := utils.InitLogger(); err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}

	log.Println("🚀 Starting Tray App...")

	// Запуск приложения
	application := app.New()
	if err := application.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
