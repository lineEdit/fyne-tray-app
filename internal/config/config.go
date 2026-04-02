package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// Config хранит настройки приложения
type Config struct {
	mu sync.RWMutex `json:"-"`

	AutoStart    bool   `json:"auto_start"`
	CheckUpdates bool   `json:"check_updates"`
	Language     string `json:"language"`
	LogLevel     string `json:"log_level"`

	WindowPosition struct {
		X, Y int `json:"x,y"`
	} `json:"window_position"`

	configPath string `json:"-"`

	// ✅ Флаги для отслеживания изменений
	loaded bool `json:"-"`
	dirty  bool `json:"-"` // Нужно ли сохранять
}

var (
	instance *Config
	once     sync.Once
)

// markDirty помечает конфиг как изменённый
func (c *Config) markDirty() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.dirty = true
}

// Get возвращает глобальный экземпляр конфигурации
func Get() *Config {
	once.Do(func() {
		instance = Load()
	})
	return instance
}

// Default возвращает конфигурацию по умолчанию
func Default() *Config {
	return &Config{
		AutoStart:    false,
		CheckUpdates: true,
		Language:     "ru-RU",
		LogLevel:     "info",
	}
}

// getConfigPath возвращает путь к файлу конфигурации рядом с exe
func getConfigPath() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	execDir := filepath.Dir(execPath)

	// ✅ Для разработки: если запускаем через go run, используем корень проекта
	// Определяем по наличию go.mod в текущей директории
	if _, err := os.Stat("go.mod"); err == nil {
		// Мы в корне проекта — сохраняем конфиг здесь
		return "config.json", nil
	}

	// Для релиза: рядом с exe
	return filepath.Join(execDir, "config.json"), nil
}

// Load загружает конфигурацию из файла рядом с exe
// Если файла нет — создаёт дефолтный
func Load() *Config {
	cfg := Default()

	configPath, err := getConfigPath()
	if err != nil {
		log.Printf("⚠️ Failed to get config path: %v", err)
		return cfg
	}
	cfg.configPath = configPath

	// Проверяем, существует ли файл
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Printf("📄 Config not found, creating default: %s", configPath)
		if saveErr := cfg.Save(); saveErr != nil {
			log.Printf("⚠️ Failed to create default config: %v", saveErr)
		}
		return cfg
	}

	// Читаем файл
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Printf("⚠️ Failed to read config: %v", err)
		return cfg
	}

	// Парсим
	if err := json.Unmarshal(data, cfg); err != nil {
		log.Printf("⚠️ Failed to parse config: %v", err)
		return cfg
	}

	log.Printf("✅ Config loaded from: %s", configPath)
	return cfg
}

// Save сохраняет только если есть изменения
func (c *Config) Save() error {
	c.mu.RLock()
	if !c.dirty && c.loaded {
		c.mu.RUnlock()
		return nil // Нет изменений — не сохраняем
	}
	// Снимаем read-lock, берём write-lock
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	// Если путь не установлен — получаем его
	if c.configPath == "" {
		path, err := getConfigPath()
		if err != nil {
			return err
		}
		c.configPath = path
	}

	// Создаём директорию
	dir := filepath.Dir(c.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Сериализуем и записываем
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(c.configPath, data, 0644); err != nil {
		return err
	}

	// ✅ Сбрасываем флаг после успешного сохранения
	c.dirty = false
	c.loaded = true

	log.Printf("💾 Config saved: %s", c.configPath)
	return nil
}

// SetLanguage с пометкой изменения
func (c *Config) SetLanguage(lang string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if lang != c.Language {
		c.Language = lang
		c.dirty = true // ✅ Помечаем как изменённый
		log.Printf("🌐 Language changed: %s", lang)
		return true
	}
	return false
}

// GetLanguage безопасно возвращает текущий язык
func (c *Config) GetLanguage() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Language
}

// Reload перечитывает конфиг с диска (если изменился извне)
func (c *Config) Reload() error {
	if c.configPath == "" {
		return nil
	}

	data, err := os.ReadFile(c.configPath)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	return json.Unmarshal(data, c)
}

// Аналогично для других полей:
func (c *Config) SetAutoStart(val bool) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if val != c.AutoStart {
		c.AutoStart = val
		c.dirty = true
		return true
	}
	return false
}

func (c *Config) SetCheckUpdates(val bool) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if val != c.CheckUpdates {
		c.CheckUpdates = val
		c.dirty = true
		return true
	}
	return false
}
