package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	AutoStart      bool   `json:"auto_start"`
	CheckUpdates   bool   `json:"check_updates"`
	Language       string `json:"language"`
	LogLevel       string `json:"log_level"`
	WindowPosition struct {
		X, Y int `json:"x,y"`
	} `json:"window_position"`
}

func Default() *Config {
	return &Config{
		AutoStart:    false,
		CheckUpdates: true,
		Language:     "ru-RU",
		LogLevel:     "info",
	}
}

func Load() *Config {
	cfg := Default()

	// Путь к конфигу в %LOCALAPPDATA%
	configDir := filepath.Join(os.Getenv("LOCALAPPDATA"), "fyne-tray-app")
	configPath := filepath.Join(configDir, "config.json")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return cfg
	}

	_ = json.Unmarshal(data, cfg)
	return cfg
}

func (c *Config) Save() error {
	dir := filepath.Join(os.Getenv("LOCALAPPDATA"), "fyne-tray-app")
	_ = os.MkdirAll(dir, 0755)

	data, _ := json.MarshalIndent(c, "", "  ")
	return os.WriteFile(filepath.Join(dir, "config.json"), data, 0644)
}
