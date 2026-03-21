package app

import (
	"log"
	"sync"

	"fyne-tray-app/internal/config"
	"fyne-tray-app/internal/tray"
	"fyne-tray-app/internal/ui"
	"fyne-tray-app/internal/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type Application struct {
	fyneApp    fyne.App
	mainWindow fyne.Window
	cfg        *config.Config
	trayMgr    *tray.Manager
}

func New() *Application {
	fyneApp := app.NewWithID("com.example.trayapp")

	// Инициализация локализации
	loc := utils.GetLocale()
	if cfg := config.Load(); cfg.Language != "" {
		_ = loc.SetLocale(cfg.Language)
	}

	return &Application{
		fyneApp: fyneApp,
		cfg:     config.Load(),
	}
}

func (a *Application) Run() error {
	log.Println("🪟 Creating main window...")

	// 1. Создаём окно
	a.mainWindow = ui.CreateMainWindow(a.fyneApp, a.cfg)

	// 2. Перехват закрытия: скрывать, а не закрывать
	a.mainWindow.SetCloseIntercept(func() {
		log.Println("🪟 Window close intercepted - hiding")
		a.mainWindow.Hide()
	})

	// 3. Инициализируем трей
	a.trayMgr = tray.NewManager(a.mainWindow, a.cfg)

	// 4. Синхронизация: ждём инициализации systray
	var ready sync.WaitGroup
	ready.Add(1)

	// 5. ✅ Запускаем systray в ГОРУТИНЕ (не блокирует главный поток)
	go func() {
		log.Println("🔌 Starting systray in goroutine...")
		err := a.trayMgr.RunWithReady(func() {
			// Callback: systray готов
			ready.Done()
			log.Println("✅ systray initialized")
		})
		if err != nil {
			log.Println(err)
			return
		}
	}()

	// 6. Ждём, пока systray инициализируется
	ready.Wait()
	log.Println("🚀 systray ready, starting Fyne event loop...")

	// 7. ✅ Запускаем Fyne event loop в ГЛАВНОМ потоке
	// Это критично для отрисовки окон!
	a.fyneApp.Run()

	return nil
}
