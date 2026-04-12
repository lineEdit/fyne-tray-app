// internal/ui/settings_window.go
package ui

import (
	"fyne-tray-app/internal/config"
	"fyne-tray-app/internal/utils"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var settingsWindowInstance fyne.Window

// ShowSettingsWindow открывает или фокусирует окно настроек
func ShowSettingsWindow(app fyne.App, cfg *config.Config, parent fyne.Window) {
	log.Println("⚙️ ShowSettingsWindow called")

	if settingsWindowInstance != nil {
		log.Println("⚙️ Settings window already exists, focusing...")
		settingsWindowInstance.Show()
		settingsWindowInstance.RequestFocus()
		return
	}

	log.Println("⚙️ Creating new settings window...")
	settingsWindowInstance = createSettingsWindow(app, cfg, parent)
	settingsWindowInstance.Show()
	settingsWindowInstance.RequestFocus()
}

func createSettingsWindow(app fyne.App, cfg *config.Config, parent fyne.Window) fyne.Window {
	loc := utils.GetLocale()

	win := app.NewWindow(loc.Get("settings.window.title"))
	win.SetContent(createSettingsContent(cfg, win))
	win.Resize(fyne.NewSize(500, 400))

	win.SetOnClosed(func() {
		log.Println("⚙️ Settings window closing, resetting instance...")
		settingsWindowInstance = nil
		if parent != nil {
			parent.RequestFocus()
		}
	})

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

	// ✅ Язык — используем language_name из файлов локалей
	availableLocales := utils.AvailableLocales()
	log.Printf("🌐 Available locales: %v", availableLocales)

	// Создаём мапу: код → отображаемое имя (для заполнения ComboBox)
	localeToName := make(map[string]string)
	for _, localeInfo := range availableLocales {
		fullInfo := utils.GetLocaleInfo(localeInfo.Code)
		localeToName[localeInfo.Code] = fullInfo.Name
		log.Printf("🌐 Locale %s: %s", localeInfo.Code, fullInfo.Name)
	}

	// Создаём список для ComboBox
	langOptions := make([]string, 0, len(localeToName))
	for _, name := range localeToName {
		langOptions = append(langOptions, name)
	}

	langSelect := widget.NewSelect(langOptions, func(value string) {
		log.Printf("🌐 ComboBox value selected: [%s]", value)

		// ✅ Надёжный поиск: сравниваем с триммингом и в нижнем регистре
		var newLang string
		selectedTrimmed := strings.TrimSpace(strings.ToLower(value))

		for code, name := range localeToName {
			nameTrimmed := strings.TrimSpace(strings.ToLower(name))
			log.Printf("🌐 Comparing: [%s] == [%s] ?", selectedTrimmed, nameTrimmed)

			if selectedTrimmed == nameTrimmed {
				newLang = code
				log.Printf("🌐 Match found: %s", newLang)
				break
			}
		}

		if newLang == "" {
			log.Printf("⚠️ No match found for: [%s]", value)
			// Fallback: ищем точное совпадение без нормализации
			for code, name := range localeToName {
				if name == value {
					newLang = code
					log.Printf("🌐 Fallback match: %s", newLang)
					break
				}
			}
		}

		log.Printf("🌐 Language selected: %s (%s)", value, newLang)

		// ✅ Сохраняем ТОЛЬКО если язык действительно изменился
		if newLang != "" && newLang != cfg.Language {
			if cfg.SetLanguage(newLang) {
				_ = utils.GetLocale().SetLocale(newLang)
				_ = cfg.Save()
				log.Printf("💾 Language saved: %s", newLang)

				dialog.ShowInformation(
					loc.Get("settings.language"),
					loc.Get("settings.language.restart_required"),
					win,
				)
			}
		} else if newLang == cfg.Language {
			log.Printf("ℹ️ Language already set to: %s", newLang)
		} else if newLang == "" {
			log.Printf("❌ Could not determine locale code")
		}
	})

	// Установка текущего значения
	currentName := localeToName[cfg.Language]
	if currentName == "" {
		// Fallback если имя не найдено
		currentName = cfg.Language
	}
	langSelect.SetSelected(currentName)
	log.Printf("🌐 ComboBox initialized with: %s (%s)", cfg.Language, currentName)

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
					logLevelSelect.SetSelected(cfg.LogLevel)

					// Обновляем выбор языка
					currentLangInfo := utils.GetLocaleInfo(cfg.Language)
					langSelect.SetSelected(currentLangInfo.Name)

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
