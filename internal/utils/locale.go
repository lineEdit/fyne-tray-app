// internal/utils/locale.go

package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// LocaleInfo — информация о доступной локали
type LocaleInfo struct {
	Code string `json:"-"`             // Код локали (из имени файла)
	Name string `json:"language_name"` // Отображаемое название
}

type LocaleManager struct {
	mu       sync.RWMutex
	locale   string
	messages map[string]string
	fallback map[string]string
}

var localeInstance *LocaleManager
var localeOnce sync.Once

func GetLocale() *LocaleManager {
	localeOnce.Do(func() {
		localeInstance = &LocaleManager{
			locale:   "ru-RU",
			messages: make(map[string]string),
			fallback: make(map[string]string),
		}
		_ = localeInstance.loadLocale("en-US")
		localeInstance.fallback = make(map[string]string)
		for k, v := range localeInstance.messages {
			localeInstance.fallback[k] = v
		}
		_ = localeInstance.loadLocale("ru-RU")
	})
	return localeInstance
}

func getLocalePath(locale string) (string, error) {
	if execPath, err := os.Executable(); err == nil {
		execDir := filepath.Dir(execPath)
		path := filepath.Join(execDir, "resources", "locales", fmt.Sprintf("%s.json", locale))
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	path := filepath.Join("resources", "locales", fmt.Sprintf("%s.json", locale))
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}
	if wd, err := os.Getwd(); err == nil {
		path := filepath.Join(wd, "resources", "locales", fmt.Sprintf("%s.json", locale))
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("locale file not found: %s.json", locale)
}

func (lm *LocaleManager) loadLocale(locale string) error {
	log.Printf("🌐 Loading locale: %s", locale) // ← Добавьте лог

	lm.mu.Lock()
	defer lm.mu.Unlock()

	path, err := getLocalePath(locale)
	if err != nil {
		log.Printf("❌ Locale path error for %s: %v", locale, err) // ← Лог ошибки
		return err
	}
	log.Printf("📂 Locale path: %s", path) // ← Путь к файлу

	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("❌ Failed to read locale file %s: %v", path, err) // ← Лог ошибки
		return fmt.Errorf("failed to read locale %s: %w", locale, err)
	}

	var messages map[string]string
	if err := json.Unmarshal(data, &messages); err != nil {
		log.Printf("❌ Failed to parse locale %s: %v", locale, err) // ← Лог ошибки
		return fmt.Errorf("failed to parse locale %s: %w", locale, err)
	}

	log.Printf("✅ Locale loaded: %s (%d keys)", locale, len(messages)) // ← Успех
	lm.messages = messages
	lm.locale = locale
	return nil
}

func (lm *LocaleManager) SetLocale(locale string) error {
	return lm.loadLocale(locale)
}

func (lm *LocaleManager) Get(key string) string {
	lm.mu.RLock()
	defer lm.mu.RUnlock()

	if msg, ok := lm.messages[key]; ok {
		return msg
	}
	if msg, ok := lm.fallback[key]; ok {
		return msg
	}
	return key
}

func (lm *LocaleManager) GetLocale() string {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	return lm.locale
}

// ✅ AvailableLocales — сканирует папку и возвращает список локалей с названиями
func AvailableLocales() []LocaleInfo {
	locales := []LocaleInfo{}

	// Пути для поиска
	searchPaths := []string{
		"resources/locales",
	}

	// Добавляем путь рядом с exe для релиза
	if execPath, err := os.Executable(); err == nil {
		execDir := filepath.Dir(execPath)
		searchPaths = append(searchPaths, filepath.Join(execDir, "resources", "locales"))
	}

	// Ищем в каждом пути
	for _, basePath := range searchPaths {
		entries, err := os.ReadDir(basePath)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			name := entry.Name()
			if strings.HasSuffix(name, ".json") {
				localeCode := strings.TrimSuffix(name, ".json")

				if isValidLocaleCode(localeCode) {
					// Проверяем, нет ли уже в списке
					exists := false
					for _, loc := range locales {
						if loc.Code == localeCode {
							exists = true
							break
						}
					}

					if !exists {
						// Читаем файл, чтобы получить language_name
						filePath := filepath.Join(basePath, name)
						data, err := os.ReadFile(filePath)
						if err == nil {
							var localeData map[string]interface{}
							if err := json.Unmarshal(data, &localeData); err == nil {
								localeName := ""
								if nameVal, ok := localeData["language_name"]; ok {
									if strVal, ok := nameVal.(string); ok {
										localeName = strVal
									}
								}

								// Если language_name не найден — используем код локали
								if localeName == "" {
									localeName = localeCode
								}

								locales = append(locales, LocaleInfo{
									Code: localeCode,
									Name: localeName,
								})
								log.Printf("🌐 Found locale: %s (%s)", localeCode, localeName)
							}
						}
					}
				}
			}
		}

		// Если нашли локали — не продолжаем поиск
		if len(locales) > 0 {
			break
		}
	}

	// Fallback: если ничего не найдено
	if len(locales) == 0 {
		log.Println("⚠️ No locales found, using defaults")
		locales = []LocaleInfo{
			{Code: "ru-RU", Name: "Русский"},
			{Code: "en-US", Name: "English"},
		}
	}

	return locales
}

func isValidLocaleCode(code string) bool {
	if len(code) < 2 {
		return false
	}
	parts := strings.Split(code, "-")
	if len(parts) == 1 {
		return len(parts[0]) == 2
	}
	if len(parts) == 2 {
		return len(parts[0]) == 2 && len(parts[1]) == 2
	}
	return false
}
