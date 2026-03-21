package widgets

import (
	"fyne-tray-app/internal/config"
	"fyne-tray-app/internal/utils"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type StatusWidget struct {
	widget.BaseWidget
	cfg   *config.Config
	label *widget.Label
	icon  *canvas.Circle
}

func NewStatusWidget(cfg *config.Config) *StatusWidget {
	w := &StatusWidget{
		cfg:   cfg,
		label: widget.NewLabel("Статус: Активно"),
		icon:  canvas.NewCircle(color.RGBA{R: 76, G: 175, B: 80, A: 255}), // Зелёный
	}
	w.ExtendBaseWidget(w)
	return w
}

func (w *StatusWidget) CreateRenderer() fyne.WidgetRenderer {
	loc := utils.GetLocale()

	content := container.NewHBox(
		w.icon,
		container.NewVBox(
			widget.NewLabelWithStyle(loc.Get("status.label"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			w.label,
		),
	)
	return widget.NewSimpleRenderer(content)
}

func (w *StatusWidget) SetStatus(status string, active bool) {
	loc := utils.GetLocale()

	if status == "" {
		if active {
			status = loc.Get("status.active")
		} else {
			status = loc.Get("status.inactive")
		}
	}
	w.label.SetText(loc.Get("status.label") + ": " + status)

	if active {
		w.icon.FillColor = color.RGBA{R: 76, G: 175, B: 80, A: 255} // Зелёный
	} else {
		w.icon.FillColor = color.RGBA{R: 255, G: 152, B: 0, A: 255} // Оранжевый
	}
	w.icon.Refresh()
	w.Refresh()
}
