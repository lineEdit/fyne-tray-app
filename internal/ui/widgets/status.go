// internal/ui/widgets/status.go
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
	loc := utils.GetLocale()

	w := &StatusWidget{
		cfg:   cfg,
		label: widget.NewLabel(loc.Get("status.label") + ": " + loc.Get("status.active")),
		icon:  canvas.NewCircle(color.RGBA{R: 76, G: 175, B: 80, A: 255}),
	}
	w.ExtendBaseWidget(w)

	// ✅ НЕ запускайте горутины здесь!
	// Если нужно обновление — делайте это через Refresh() по событию

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

// SetStatus обновляет статус (вызывайте из основного потока Fyne)
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
		w.icon.FillColor = color.RGBA{R: 76, G: 175, B: 80, A: 255}
	} else {
		w.icon.FillColor = color.RGBA{R: 255, G: 152, B: 0, A: 255}
	}
	w.icon.Refresh()
	w.Refresh()
}
