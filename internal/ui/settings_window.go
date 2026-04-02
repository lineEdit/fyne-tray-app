// internal/ui/settings_window.go

package ui

import (
	"fyne-tray-app/internal/config"
	"fyne-tray-app/internal/utils"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var settingsWindowInstance fyne.Window

func ShowSettingsWindow(app fyne.App, cfg *config.Config, parent fyne.Window) {
	if settingsWindowInstance != nil {
		settingsWindowInstance.Show()
		settingsWindowInstance.RequestFocus()
		return
	}

	settingsWindowInstance = createSettingsWindow(app, cfg, parent)
	settingsWindowInstance.Show()
	settingsWindowInstance.RequestFocus()
}

func createSettingsWindow(app fyne.App, cfg *config.Config, parent fyne.Window) fyne.Window {
	loc := utils.GetLocale()

	win := app.NewWindow(loc.Get("settings.window.title"))
	win.SetContent(createSettingsContent(cfg, win))
	win.Resize(fyne.NewSize(500, 400))

	if parent != nil {
		win.SetOnClosed(func() {
			parent.RequestFocus()
		})
	}

	log.Println("⚙️ Settings window created")
	return win
}

func createSettingsContent(cfg *config.Config, win fyne.Window) fyne.CanvasObject {
	loc := utils.GetLocale()

	// Автозапуск
	autoStart := widget.NewCheck(loc.Get("settings.auto_start"), func(checked bool) {
		if cfg.SetAutoStart(checked) {
			_ = cfg.Save()
		}
	})
	autoStart.Checked = cfg.AutoStart

	// Проверка обновлений
	checkUpdates := widget.NewCheck(loc.Get("settings.check_updates"), func(checked bool) {
		if cfg.SetCheckUpdates(checked) {
			_ = cfg.Save()
		}
	})
	checkUpdates.Checked = cfg.CheckUpdates

	// ✅ Язык — динамический список из AvailableLocales()
	availableLocales := utils.AvailableLocales()
	log.Printf("🌐 Available locales: %v", availableLocales)

	// Создаём список названий для отображения
	langOptions := make([]string, len(availableLocales))
	for i, localeInfo := range availableLocales {
		langOptions[i] = localeInfo.Name
	}

	langSelect := widget.NewSelect(langOptions, func(value string) {
		// Находим код локали по названию
		var newLang string
		for _, localeInfo := range availableLocales {
			if localeInfo.Name == value {
				newLang = localeInfo.Code
				break
			}
		}

		log.Printf("🌐 Language selected: %s (%s)", value, newLang)

		if newLang != "" && cfg.SetLanguage(newLang) {
			_ = utils.GetLocale().SetLocale(newLang)
			_ = cfg.Save()

			dialog.ShowInformation(
				loc.Get("settings.language"),
				loc.Get("settings.language.restart_required"),
				win,
			)
		}
	})

	// Установка текущего значения
	for _, localeInfo := range availableLocales {
		if localeInfo.Code == cfg.Language {
			langSelect.SetSelected(localeInfo.Name)
			break
		}
	}

	// Уровень логирования
	logLevelSelect := widget.NewSelect([]string{"debug", "info", "warn", "error"}, func(value string) {
		if cfg.SetLogLevel(value) {
			_ = cfg.Save()
		}
	})
	logLevelSelect.SetSelected(cfg.LogLevel)

	// Кнопки
	saveBtn := widget.NewButton(loc.Get("settings.btn.save"), func() {
		_ = cfg.Save()
		dialog.ShowInformation(
			loc.Get("settings.notification.title"),
			loc.Get("settings.notification.saved"),
			win,
		)
	})

	resetBtn := widget.NewButton(loc.Get("settings.btn.reset"), func() {
		dialog.ShowConfirm(
			loc.Get("settings.btn.reset"),
			loc.Get("settings.reset.confirm"),
			func(ok bool) {
				if ok {
					defaultCfg := config.Default()
					cfg.AutoStart = defaultCfg.AutoStart
					cfg.CheckUpdates = defaultCfg.CheckUpdates
					cfg.Language = defaultCfg.Language
					cfg.LogLevel = defaultCfg.LogLevel
					_ = cfg.Save()

					autoStart.Checked = cfg.AutoStart
					autoStart.Refresh()
					checkUpdates.Checked = cfg.CheckUpdates
					checkUpdates.Refresh()

					// Обновляем выбор языка
					for _, localeInfo := range availableLocales {
						if localeInfo.Code == cfg.Language {
							langSelect.SetSelected(localeInfo.Name)
							break
						}
					}

					logLevelSelect.SetSelected(cfg.LogLevel)

					dialog.ShowInformation(
						loc.Get("settings.notification.title"),
						loc.Get("settings.notification.reset"),
						win,
					)
				}
			},
			win,
		)
	})

	closeBtn := widget.NewButton(loc.Get("settings.btn.close"), func() {
		win.Close()
	})

	forms := container.NewVBox(
		widget.NewForm(
			widget.NewFormItem(loc.Get("settings.auto_start"), autoStart),
			widget.NewFormItem(loc.Get("settings.check_updates"), checkUpdates),
			widget.NewFormItem(loc.Get("settings.language"), langSelect),
			widget.NewFormItem(loc.Get("settings.log_level"), logLevelSelect),
		),
		widget.NewSeparator(),
		container.NewHBox(saveBtn, resetBtn),
		container.NewCenter(closeBtn),
	)

	return container.NewPadded(forms)
}
