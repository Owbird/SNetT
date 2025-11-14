package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/owbird/snett/internal/ui"
)

func main() {
	a := app.New()

	w := a.NewWindow("SNetT")

	w.Resize(fyne.NewSize(500, 500))

	hui := ui.NewHomeUI(w)
	sui := ui.NewServerUI(w)
	wui := ui.NewWormholeUI(w)

	menus := []*fyne.Menu{
		fyne.NewMenu("Server",
			fyne.NewMenuItem("Host Directory", sui.ChooseHostDir),
			fyne.NewMenuItem("Discover", sui.Discover),
			fyne.NewMenuItem("Settings", sui.ServerSettings),
		),
		fyne.NewMenu("Wormhole",
			fyne.NewMenuItem("Share file", wui.ShareFile),
			fyne.NewMenuItem("Receive file", wui.ReceiveFile),
		),

	}

	w.SetMainMenu(fyne.NewMainMenu(menus...))

	w.SetContent(hui.Home())

	w.Show()

	a.Run()
}
