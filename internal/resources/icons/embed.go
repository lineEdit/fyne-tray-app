// Package icons предоставляет встроенные ресурсы.
package icons

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed icon-on.ico
var iconData []byte

// ResourceIconOnIco — экспортированная иконка для системного трея
var ResourceIconOnIco = &fyne.StaticResource{
	StaticName:    "icon-on.ico",
	StaticContent: iconData,
}
