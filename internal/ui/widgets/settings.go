package widgets

import (
	"fyne-tray-app/internal/config"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type SettingsWidget struct {
	widget.BaseWidget
	cfg       *config.Config
	expanded  bool
	toggleBtn *widget.Button
	content   *fyne.Container
}

func NewSettingsWidget(cfg *config.Config) *SettingsWidget {
	w := &SettingsWidget{
		cfg:      cfg,
		expanded: false,
	}
	w.ExtendBaseWidget(w)
	return w
}

func (w *SettingsWidget) CreateRenderer() fyne.WidgetRenderer {
	w.toggleBtn = widget.NewButton("Настройки ⌄", w.toggleExpand)
	w.content = w.createSettingsContent()
	w.content.Hidden = true

	contents := container.NewVBox(w.toggleBtn, w.content)
	return widget.NewSimpleRenderer(contents)
}

func (w *SettingsWidget) toggleExpand() {
	w.expanded = !w.expanded
	if w.expanded {
		w.toggleBtn.SetText("Настройки ⌃")
		w.content.Hidden = false
	} else {
		w.toggleBtn.SetText("Настройки ⌄")
		w.content.Hidden = true
	}
	w.content.Refresh()
	w.toggleBtn.Refresh()
}

func (w *SettingsWidget) Expand() {
	if !w.expanded {
		w.toggleExpand()
	}
}

func (w *SettingsWidget) createSettingsContent() *fyne.Container {
	autoStart := widget.NewCheck("Автозапуск", func(checked bool) {
		w.cfg.AutoStart = checked
		_ = w.cfg.Save()
	})
	autoStart.Checked = w.cfg.AutoStart

	checkUpdates := widget.NewCheck("Проверять обновления", func(checked bool) {
		w.cfg.CheckUpdates = checked
		_ = w.cfg.Save()
	})
	checkUpdates.Checked = w.cfg.CheckUpdates

	lang := widget.NewSelect([]string{"Русский", "English"}, func(value string) {
		if value == "Русский" {
			w.cfg.Language = "ru-RU"
		} else {
			w.cfg.Language = "en-US"
		}
		_ = w.cfg.Save()
	})
	if w.cfg.Language == "ru-RU" {
		lang.SetSelected("Русский")
	} else {
		lang.SetSelected("English")
	}

	saveBtn := widget.NewButton("Применить", func() {
		_ = w.cfg.Save()
		fyne.CurrentApp().SendNotification(&fyne.Notification{
			Title:   "Настройки",
			Content: "Сохранено",
		})
	})

	return container.NewPadded(container.NewVBox(
		autoStart,
		checkUpdates,
		container.NewHBox(widget.NewLabel("Язык:"), lang),
		container.NewCenter(saveBtn),
	))
}
