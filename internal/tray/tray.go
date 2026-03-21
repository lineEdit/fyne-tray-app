package tray

import (
	"fyne-tray-app/internal/config"
	"fyne-tray-app/internal/resources/icons"
	"fyne-tray-app/internal/utils"
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

// RunWithReady запускает systray с callback при готовности
func (m *Manager) RunWithReady(onReady func()) error {
	systray.Run(m.createReadyHandler(onReady), m.exit)
	return nil
}

// Run — обратная совместимость
func (m *Manager) Run() error {
	return m.RunWithReady(nil)
}

// createReadyHandler создаёт обработчик onReady с поддержкой синхронизации
func (m *Manager) createReadyHandler(onReady func()) func() {
	return func() {
		loc := utils.GetLocale()

		systray.SetTitle(loc.Get("tray.title"))
		systray.SetTooltip(loc.Get("tray.tooltip"))

		// Иконка
		if icons.ResourceIconOnIco != nil {
			systray.SetIcon(icons.ResourceIconOnIco.Content())
			log.Println("✅ Icon loaded from resources")
		}

		// Меню
		mShow := systray.AddMenuItem(loc.Get("tray.menu.show"), "")
		mQuit := systray.AddMenuItem(loc.Get("tray.menu.exit"), "")

		// ✅ Сигнализируем готовность (если передан callback)
		if onReady != nil {
			onReady()
		}

		// ✅ Обработчики кликов в отдельной горутине
		go func() {
			for {
				select {
				case <-mShow.ClickedCh:
					log.Println("🔘 Show menu clicked")
					// ✅ UI-операции только через fyne.Do
					fyne.Do(func() {
						if m.window == nil {
							log.Println("❌ Window is nil!")
							return
						}
						log.Println("🪟 Showing window...")
						m.window.Show()
						m.window.RequestFocus()
						m.window.Resize(fyne.NewSize(600, 400))
					})
				case <-mQuit.ClickedCh:
					log.Println("👋 Quit menu clicked")
					systray.Quit()
					return
				}
			}
		}()
	}
}

func (m *Manager) exit() {
	log.Println("🔚 systray exit callback")
}
