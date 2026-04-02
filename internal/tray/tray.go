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
	window       fyne.Window
	cfg          *config.Config
	menuShow     *systray.MenuItem
	menuSettings *systray.MenuItem // ← Новый пункт меню
	menuQuit     *systray.MenuItem
	onSettings   func() // ← Callback для открытия окна настроек
}

func NewManager(w fyne.Window, cfg *config.Config) *Manager {
	log.Println("🔌 Tray.NewManager called")
	return &Manager{window: w, cfg: cfg}
}

// SetOnSettingsCallback устанавливает callback для открытия окна настроек
func (m *Manager) SetOnSettingsCallback(callback func()) {
	m.onSettings = callback
}

// RunWithReady запускает systray с callback при готовности
func (m *Manager) RunWithReady(onReady func()) error {
	log.Println("🔌 Tray.RunWithReady: entering systray.Run")

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
	loc := utils.GetLocale()

	log.Println("🪟 tray.ready: setting title/tooltip")
	systray.SetTitle(loc.Get("tray.title"))
	systray.SetTooltip(loc.Get("tray.tooltip"))

	// Иконка
	log.Println("🖼️ tray.ready: loading icon...")
	if icons.ResourceIconOnIco != nil {
		log.Printf("🖼️ Icon resource found: %d bytes", len(icons.ResourceIconOnIco.Content()))
		systray.SetIcon(icons.ResourceIconOnIco.Content())
		log.Println("✅ Icon set in tray")
	}

	// Меню
	log.Println("📋 tray.ready: creating menu items")
	m.menuShow = systray.AddMenuItem(loc.Get("tray.menu.show"), "")

	// ✅ Новый пункт: Настройки
	m.menuSettings = systray.AddMenuItem(loc.Get("settings.menu"), "")

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

			// ✅ Обработчик настроек
			case <-m.menuSettings.ClickedCh:
				log.Println("⚙️ Settings clicked")
				fyne.Do(func() {
					if m.onSettings != nil {
						log.Println("⚙️ Opening settings window...")
						m.onSettings()
					} else {
						log.Println("⚠️ Settings callback not set")
					}
				})

			case <-m.menuQuit.ClickedCh:
				log.Println("👋 Quit clicked")
				fyne.Do(func() {
					log.Println("🪟 Closing window and Fyne app...")
					if m.window != nil {
						m.window.Close()
					}
					fyne.CurrentApp().Quit()
				})
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
