package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// LocaleManager управляет локализацией приложения
type LocaleManager struct {
	mu       sync.RWMutex
	locale   string
	messages map[string]string
	fallback map[string]string
}

var localeInstance *LocaleManager
var localeOnce sync.Once

// GetLocale возвращает глобальный экземпляр LocaleManager
func GetLocale() *LocaleManager {
	localeOnce.Do(func() {
		localeInstance = &LocaleManager{
			locale:   "ru-RU",
			messages: make(map[string]string),
			fallback: make(map[string]string),
		}
		// Загружаем fallback (английский)
		_ = localeInstance.loadLocale("en-US")
		localeInstance.fallback = make(map[string]string)
		for k, v := range localeInstance.messages {
			localeInstance.fallback[k] = v
		}
		// Загружаем основную локаль
		_ = localeInstance.loadLocale("ru-RU")
	})
	return localeInstance
}

// getLocalePath возвращает путь к файлу локализации
func getLocalePath(locale string) (string, error) {
	// 1. Пробуем путь относительно exe (для релиза)
	if execPath, err := os.Executable(); err == nil {
		execDir := filepath.Dir(execPath)
		path := filepath.Join(execDir, "resources", "locales", fmt.Sprintf("%s.json", locale))
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	// 2. Пробуем путь относительно рабочей директории (для go run)
	path := filepath.Join("resources", "locales", fmt.Sprintf("%s.json", locale))
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}

	// 3. Пробуем от корня проекта (для разработки в IDE)
	if wd, err := os.Getwd(); err == nil {
		path := filepath.Join(wd, "resources", "locales", fmt.Sprintf("%s.json", locale))
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("locale file not found: %s.json", locale)
}

// loadLocale загружает файл локализации из файловой системы
func (lm *LocaleManager) loadLocale(locale string) error {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	path, err := getLocalePath(locale)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read locale %s: %w", locale, err)
	}

	var messages map[string]string
	if err := json.Unmarshal(data, &messages); err != nil {
		return fmt.Errorf("failed to parse locale %s: %w", locale, err)
	}

	lm.messages = messages
	lm.locale = locale
	return nil
}

// SetLocale устанавливает текущую локаль
func (lm *LocaleManager) SetLocale(locale string) error {
	return lm.loadLocale(locale)
}

// Get возвращает строку для текущего языка
func (lm *LocaleManager) Get(key string) string {
	lm.mu.RLock()
	defer lm.mu.RUnlock()

	if msg, ok := lm.messages[key]; ok {
		return msg
	}
	// Fallback на английский
	if msg, ok := lm.fallback[key]; ok {
		return msg
	}
	// Если ничего не найдено — возвращаем ключ
	return key
}

// GetLocale возвращает текущую локаль
func (lm *LocaleManager) GetLocale() string {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	return lm.locale
}

// AvailableLocales возвращает список доступных локалей
func AvailableLocales() []string {
	return []string{"ru-RU", "en-US"}
}
