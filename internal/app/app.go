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
	// 1. Создаём приложение Fyne
	fyneApp := app.NewWithID("com.example.trayapp")

	// 2. ✅ Загружаем конфиг ПЕРВЫМ (создаст дефолтный рядом с exe, если нет)
	// Используем config.Load() напрямую, а не синглтон, для явного контроля
	cfg := config.Load()

	// 3. ✅ Применяем язык из конфига к локализации ДО создания окна
	loc := utils.GetLocale()
	if cfg.Language != "" {
		if err := loc.SetLocale(cfg.Language); err != nil {
			log.Printf("⚠️ Failed to set locale %s: %v", cfg.Language, err)
		} else {
			log.Printf("🌐 Locale applied: %s", cfg.Language)
		}
	}

	return &Application{
		fyneApp: fyneApp,
		cfg:     cfg, // Сохраняем ссылку на тот же экземпляр конфига
	}
}

func (a *Application) Run() error {
	log.Println("🪟 [1/7] Creating main window...")
	a.mainWindow = ui.CreateMainWindow(a.fyneApp, a.cfg)
	log.Println("🪟 [2/7] Window created")

	log.Println("🔒 [3/7] Setting close intercept...")
	a.mainWindow.SetCloseIntercept(func() {
		log.Println("🪟 Window close intercepted - hiding")
		a.mainWindow.Hide()
	})
	log.Println("🔒 [4/7] Close intercept set")

	log.Println("🔌 [5/7] Creating tray manager...")
	a.trayMgr = tray.NewManager(a.mainWindow, a.cfg)

	// ✅ Устанавливаем callback для открытия окна настроек
	a.trayMgr.SetOnSettingsCallback(func() {
		log.Println("⚙️ Opening settings window from tray")
		ui.ShowSettingsWindow(a.fyneApp, a.cfg, a.mainWindow)
	})

	log.Println("🔌 [6/7] Tray manager created")

	var ready sync.WaitGroup
	ready.Add(1)

	log.Println("🚀 [7/7] Starting systray goroutine...")
	go func() {
		log.Println("🔌 Goroutine: calling RunWithReady...")
		err := a.trayMgr.RunWithReady(func() {
			log.Println("✅ systray ready callback fired")
			ready.Done()
		})
		if err != nil {
			return
		}
	}()

	log.Println("⏳ Waiting for systray ready...")
	ready.Wait()
	log.Println("🚀 systray ready, starting Fyne event loop...")

	a.fyneApp.Run()

	return nil
}

// GetConfig возвращает текущую конфигурацию (для внешнего доступа)
func (a *Application) GetConfig() *config.Config {
	return a.cfg
}

// ReloadConfig перечитывает конфиг с диска (если изменился извне)
func (a *Application) ReloadConfig() error {
	if a.cfg == nil {
		return nil
	}

	oldLang := a.cfg.GetLanguage()

	if err := a.cfg.Reload(); err != nil {
		return err
	}

	// Если язык изменился — применяем к локализации
	newLang := a.cfg.GetLanguage()
	if newLang != "" && newLang != oldLang {
		log.Printf("🌐 Language changed via config reload: %s → %s", oldLang, newLang)
		if err := utils.GetLocale().SetLocale(newLang); err != nil {
			log.Printf("⚠️ Failed to apply new locale: %v", err)
		}
	}

	return nil
}
