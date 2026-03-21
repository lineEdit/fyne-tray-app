package ui

import (
	"fyne-tray-app/internal/config"
	"fyne-tray-app/internal/ui/widgets"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func CreateMainWindow(app fyne.App, cfg *config.Config) fyne.Window {
	window := app.NewWindow("Tray App")

	content := createMainContent(cfg)
	window.SetContent(content)
	window.Resize(fyne.NewSize(450, 320))
	//window.SetMinSize(fyne.NewSize(400, 280))

	return window
}

func createMainContent(cfg *config.Config) fyne.CanvasObject {
	header := container.NewHBox(
		widget.NewLabel("Tray App"),
	)

	statusWidget := widgets.NewStatusWidget(cfg)
	settingsWidget := widgets.NewSettingsWidget(cfg)

	actions := container.NewHBox(
		widget.NewButton("Скрыть", func() {
			windows := fyne.CurrentApp().Driver().AllWindows()
			if len(windows) > 0 {
				windows[0].Hide()
			}
		}),
		widget.NewButton("Настройки", func() {
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
