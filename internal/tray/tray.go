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
	window   fyne.Window
	cfg      *config.Config
	menuShow *systray.MenuItem
	menuQuit *systray.MenuItem
}

func NewManager(w fyne.Window, cfg *config.Config) *Manager {
	log.Println("🔌 tray.NewManager: ENTERED") // ← Должно появиться!
	return &Manager{window: w, cfg: cfg}
}

func (m *Manager) RunWithReady(onReady func()) error {
	log.Println("🔌 tray.RunWithReady: ENTERED") // ← Должно появиться!

	systray.Run(func() {
		log.Println("✅ systray.onReady: STARTED")
		m.ready(onReady)
		log.Println("✅ systray.onReady: COMPLETED")
	}, func() {
		log.Println("🔚 systray.onExit: called")
		m.exit()
	})

	log.Println("⚠️ systray.Run returned")
	return nil
}

func (m *Manager) ready(onReady func()) {
	log.Println("🪟 tray.ready: setting title/tooltip")

	loc := utils.GetLocale()
	systray.SetTitle(loc.Get("tray.title"))
	systray.SetTooltip(loc.Get("tray.tooltip"))

	// Иконка — с подробной отладкой
	log.Println("🖼️ tray.ready: loading icon...")
	if icons.ResourceIconOnIco != nil {
		log.Printf("🖼️ Icon resource found: %d bytes", len(icons.ResourceIconOnIco.Content()))
		systray.SetIcon(icons.ResourceIconOnIco.Content())
		log.Println("✅ Icon set in tray")
	} else {
		log.Println("❌ Icon resource is NIL — using fallback")
		// Попробуем загрузить из файла как запасной вариант
		// (необязательно, но поможет при отладке)
	}

	// Меню
	log.Println("📋 tray.ready: creating menu items")
	m.menuShow = systray.AddMenuItem(loc.Get("tray.menu.show"), "")
	m.menuQuit = systray.AddMenuItem(loc.Get("tray.menu.exit"), "")
	log.Println("✅ Menu items created")

	// ✅ Сигнализируем готовность
	if onReady != nil {
		log.Println("🔔 Calling onReady callback")
		onReady()
	}

	// Обработчики
	log.Println("🔄 tray.ready: starting event loop goroutine")
	go func() {
		for {
			select {
			case <-m.menuShow.ClickedCh:
				log.Println("🔘 Show clicked")
				fyne.Do(func() {
					if m.window == nil {
						log.Println("❌ Window is nil!")
						return
					}
					log.Println("🪟 Showing window from tray")
					m.window.Show()
					m.window.RequestFocus()
				})
			case <-m.menuQuit.ClickedCh:
				log.Println("👋 Quit clicked")

				// ✅ Оборачиваем все Fyne-операции в fyne.Do()
				fyne.Do(func() {
					// 1. Закрываем окно (если открыто)
					if m.window != nil {
						m.window.Close()
					}

					// 2. Закрываем Fyne приложение
					fyne.CurrentApp().Quit()
				})

				// 3. Закрываем трей (можно вне fyne.Do, это systray API)
				systray.Quit()
				return
			}
		}
	}()
	log.Println("✅ Event loop goroutine started")
}

func (m *Manager) exit() {
	log.Println("🔚 tray.exit: cleaning up")
}
