package app

import (
	"fyne-tray-app/internal/config"
	"fyne-tray-app/internal/tray"
	"fyne-tray-app/internal/ui"
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

	return &Application{
		fyneApp: fyneApp,
		cfg:     config.Load(),
	}
}

func (a *Application) Run() error {
	// Создаём окно через UI-модуль
	a.mainWindow = ui.CreateMainWindow(a.fyneApp, a.cfg)

	// Перехват закрытия окна
	a.mainWindow.SetCloseIntercept(func() {
		a.mainWindow.Hide()
	})

	// Инициализируем трей
	a.trayMgr = tray.NewManager(a.mainWindow, a.cfg)

	// Запускаем systray (блокирующий вызов)
	return a.trayMgr.Run()
}
