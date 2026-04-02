// internal/ui/widgets/settings.go
package widgets

import (
	"fyne-tray-app/internal/config"
	"fyne-tray-app/internal/utils"

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
	return w // ✅ Никаких горутин здесь!
}

func (w *SettingsWidget) CreateRenderer() fyne.WidgetRenderer {
	loc := utils.GetLocale()

	w.toggleBtn = widget.NewButton(loc.Get("settings.title")+" ⌄", w.toggleExpand)
	w.content = w.createSettingsContent()
	w.content.Hidden = true

	contents := container.NewVBox(w.toggleBtn, w.content)
	return widget.NewSimpleRenderer(contents)
}

func (w *SettingsWidget) toggleExpand() {
	loc := utils.GetLocale()

	w.expanded = !w.expanded
	if w.expanded {
		w.toggleBtn.SetText(loc.Get("settings.title") + " ⌃")
		w.content.Hidden = false
	} else {
		w.toggleBtn.SetText(loc.Get("settings.title") + " ⌄")
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
	loc := utils.GetLocale()

	autoStart := widget.NewCheck(loc.Get("settings.auto_start"), func(checked bool) {
		w.cfg.SetAutoStart(checked)
		_ = w.cfg.Save()
	})
	autoStart.Checked = w.cfg.AutoStart

	checkUpdates := widget.NewCheck(loc.Get("settings.check_updates"), func(checked bool) {
		w.cfg.SetCheckUpdates(checked)
		_ = w.cfg.Save()
	})
	checkUpdates.Checked = w.cfg.CheckUpdates

	lang := widget.NewSelect([]string{
		loc.Get("settings.language.ru"),
		loc.Get("settings.language.en"),
	}, func(value string) {
		var newLang string
		if value == loc.Get("settings.language.ru") {
			newLang = "ru-RU"
		} else {
			newLang = "en-US"
		}

		if w.cfg.SetLanguage(newLang) {
			_ = utils.GetLocale().SetLocale(newLang)
			_ = w.cfg.Save()

			// Пересоздание контента
			fyne.Do(func() {
				w.content = w.createSettingsContent()
				w.content.Refresh()
				w.toggleBtn.SetText(loc.Get("settings.title") +
					map[bool]string{true: " ⌃", false: " ⌄"}[w.expanded])
				w.toggleBtn.Refresh()

				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Title:   loc.Get("settings.notification.title"),
					Content: loc.Get("settings.notification.saved"),
				})
			})
		}
	})
	if w.cfg.Language == "ru-RU" {
		lang.SetSelected(loc.Get("settings.language.ru"))
	} else {
		lang.SetSelected(loc.Get("settings.language.en"))
	}

	saveBtn := widget.NewButton(loc.Get("settings.btn.save"), func() {
		fyne.CurrentApp().SendNotification(&fyne.Notification{
			Title:   loc.Get("settings.notification.title"),
			Content: loc.Get("settings.notification.saved"),
		})
	})

	return container.NewPadded(container.NewVBox(
		autoStart,
		checkUpdates,
		container.NewHBox(widget.NewLabel(loc.Get("settings.language")+":"), lang),
		container.NewCenter(saveBtn),
	))
}
