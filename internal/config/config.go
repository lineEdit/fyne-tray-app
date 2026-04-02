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
	Mu sync.RWMutex `json:"-"`

	AutoStart    bool   `json:"auto_start"`
	CheckUpdates bool   `json:"check_updates"`
	Language     string `json:"language"`
	LogLevel     string `json:"log_level"`

	WindowPosition struct {
		X int `json:"x"`
		Y int `json:"y"`
	} `json:"window_position"`

	ConfigPath string `json:"-"`

	// ✅ Флаги для отслеживания изменений
	Loaded bool `json:"-"`
	Dirty  bool `json:"-"` // Нужно ли сохранять
}

//var (
//	instance *Config
//	once     sync.Once
//)

// markDirty помечает конфиг как изменённый
func (c *Config) markDirty() {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	c.Dirty = true
}

// Get возвращает глобальный экземпляр конфигурации
//func Get() *Config {
//	once.Do(func() {
//		instance = Load()
//	})
//	return instance
//}

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
	cfg.ConfigPath = configPath

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
	c.Mu.RLock()
	if !c.Dirty && c.Loaded {
		c.Mu.RUnlock()
		return nil // Нет изменений — не сохраняем
	}
	// Снимаем read-lock, берём write-lock
	c.Mu.RUnlock()

	c.Mu.Lock()
	defer c.Mu.Unlock()

	// Если путь не установлен — получаем его
	if c.ConfigPath == "" {
		path, err := getConfigPath()
		if err != nil {
			return err
		}
		c.ConfigPath = path
	}

	// Создаём директорию
	dir := filepath.Dir(c.ConfigPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Сериализуем и записываем
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(c.ConfigPath, data, 0644); err != nil {
		return err
	}

	// ✅ Сбрасываем флаг после успешного сохранения
	c.Dirty = false
	c.Loaded = true

	log.Printf("💾 Config saved: %s", c.ConfigPath)
	return nil
}

// SetLanguage с пометкой изменения
func (c *Config) SetLanguage(lang string) bool {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	if lang != c.Language {
		c.Language = lang
		c.Dirty = true // ✅ Помечаем как изменённый
		log.Printf("🌐 Language changed: %s", lang)
		return true
	}
	return false
}

// GetLanguage безопасно возвращает текущий язык
func (c *Config) GetLanguage() string {
	c.Mu.RLock()
	defer c.Mu.RUnlock()
	return c.Language
}

// Reload перечитывает конфиг с диска (если изменился извне)
func (c *Config) Reload() error {
	if c.ConfigPath == "" {
		return nil
	}

	data, err := os.ReadFile(c.ConfigPath)
	if err != nil {
		return err
	}

	c.Mu.Lock()
	defer c.Mu.Unlock()

	return json.Unmarshal(data, c)
}

// Аналогично для других полей:
func (c *Config) SetAutoStart(val bool) bool {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	if val != c.AutoStart {
		c.AutoStart = val
		c.Dirty = true
		return true
	}
	return false
}

func (c *Config) SetCheckUpdates(val bool) bool {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	if val != c.CheckUpdates {
		c.CheckUpdates = val
		c.Dirty = true
		return true
	}
	return false
}

// SetLogLevel устанавливает уровень логирования
func (c *Config) SetLogLevel(level string) bool {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	if validLevels[level] && level != c.LogLevel {
		c.LogLevel = level
		c.Dirty = true
		log.Printf("📝 LogLevel changed: %s", level)
		return true
	}
	return false
}

// GetLogLevel безопасно возвращает уровень логирования
func (c *Config) GetLogLevel() string {
	c.Mu.RLock()
	defer c.Mu.RUnlock()
	return c.LogLevel
}
