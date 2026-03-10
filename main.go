package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/systray"
)

var mainWindow fyne.Window

//go:generate fyne bundle -o resources.go -package main icon-on.ico

func main() {
	fyneApp := app.New()
	mainWindow = fyneApp.NewWindow("Tray App")
	mainWindow.SetContent(container.NewVBox(
		widget.NewLabel("Приложение работает!"),
		widget.NewButton("Скрыть", func() { mainWindow.Hide() }),
	))
	mainWindow.Resize(fyne.NewSize(400, 300))

	mainWindow.SetCloseIntercept(func() {
		mainWindow.Hide()
	})

	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle("Tray App")
	systray.SetTooltip("Кликните для меню")

	// Загрузка иконки с отладкой
	if resourceIconOnIco != nil {
		systray.SetIcon(resourceIconOnIco.Content())
		log.Printf("✅ Иконка загружена из ресурсов: %d байт", len(resourceIconOnIco.Content()))
	}

	mShow := systray.AddMenuItem("Показать", "")
	mQuit := systray.AddMenuItem("Выход", "")

	go func() {
		for {
			select {
			case <-mShow.ClickedCh:
				fyne.Do(func() {
					mainWindow.Show()
					mainWindow.RequestFocus()
				})
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {
	log.Println("Выход из приложения")
}
