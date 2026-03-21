package tray

import (
	"fyne-tray-app/internal/config"
	"fyne-tray-app/internal/resources/icons" // ← Импорт сгенерированных ресурсов
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/systray"
)

type Manager struct {
	window fyne.Window
	cfg    *config.Config
}

func NewManager(w fyne.Window, cfg *config.Config) *Manager {
	return &Manager{window: w, cfg: cfg}
}

func (m *Manager) Run() error {
	systray.Run(m.ready, m.exit)
	return nil
}

func (m *Manager) ready() {
	systray.SetTitle("Tray App")
	systray.SetTooltip("Кликните для меню")

	// ✅ Используем сгенерированную иконку
	if icons.ResourceIconOnIco != nil {
		systray.SetIcon(icons.ResourceIconOnIco.Content())
		log.Println("✅ Иконка загружена из ресурсов")
	} else {
		log.Println("⚠️ Иконка не найдена в ресурсах")
	}

	// Меню
	mShow := systray.AddMenuItem("Показать окно", "")
	mQuit := systray.AddMenuItem("Выход", "Закрыть приложение")

	go func() {
		for {
			select {
			case <-mShow.ClickedCh:
				fyne.Do(func() {
					m.window.Show()
					m.window.RequestFocus()
				})
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func (m *Manager) exit() {
	log.Println("Выход из приложения")
}
