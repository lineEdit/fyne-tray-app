package utils

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

var Logger *log.Logger

func InitLogger() error {
	logDir := filepath.Join(os.Getenv("LOCALAPPDATA"), "fyne-tray-app", "logs")
	_ = os.MkdirAll(logDir, 0755)

	logPath := filepath.Join(logDir, "app.log")
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	mw := io.MultiWriter(file, os.Stdout)
	Logger = log.New(mw, "", log.Ldate|log.Ltime|log.Lshortfile)

	return nil
}
