package ui

import (
	"fyne-tray-app/internal/config"
	"fyne-tray-app/internal/ui/widgets"
	"fyne-tray-app/internal/utils"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func CreateMainWindow(app fyne.App, cfg *config.Config) fyne.Window {
	loc := utils.GetLocale()
	log.Println("🪟 Creating window:", loc.Get("window.title"))

	window := app.NewWindow(loc.Get("window.title"))
	content := createMainContent(cfg)
	window.SetContent(content)
	window.Resize(fyne.NewSize(450, 320))
	return window
}

func createMainContent(cfg *config.Config) fyne.CanvasObject {
	loc := utils.GetLocale()

	header := container.NewHBox(
		widget.NewLabel(loc.Get("app.name")),
	)

	statusWidget := widgets.NewStatusWidget(cfg)
	settingsWidget := widgets.NewSettingsWidget(cfg)

	actions := container.NewHBox(
		widget.NewButton(loc.Get("window.btn.hide"), func() {
			windows := fyne.CurrentApp().Driver().AllWindows()
			if len(windows) > 0 {
				windows[0].Hide()
			}
		}),
		widget.NewButton(loc.Get("window.btn.settings"), func() {
			settingsWidget.Expand()
		}),
	)

	return container.NewBorder(
		header,
		container.NewHBox(
			widget.NewLabel("v1.0.0"),
			container.NewCenter(actions),
		),
		nil, nil,
		container.NewVBox(
			statusWidget,
			widget.NewSeparator(),
			settingsWidget,
		),
	)
}
