package theme

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

type AppTheme struct {
	fyne.Theme
}

func NewAppTheme() *AppTheme {
	return &AppTheme{
		Theme: theme.DarkTheme(),
	}
}

func (t *AppTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case "primary":
		return color.NRGBA{R: 0x21, G: 0x96, B: 0xf3, A: 0xff}
	case "success":
		return color.NRGBA{R: 0x4c, G: 0xaf, B: 0x50, A: 0xff}
	case "warning":
		return color.NRGBA{R: 0xff, G: 0x98, B: 0x00, A: 0xff}
	case "error":
		return color.NRGBA{R: 0xf4, G: 0x43, B: 0x36, A: 0xff}
	}
	return t.Theme.Color(name, variant)
}
