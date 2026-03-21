// Package icons предоставляет встроенные ресурсы.
package icons

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed Icon-on.ico
var iconOnIcoData []byte

// ResourceIconOnIco — экспортированная иконка для использования в других пакетах.
var ResourceIconOnIco = &fyne.StaticResource{
	StaticName:    "icon-on.ico",
	StaticContent: iconOnIcoData,
}
