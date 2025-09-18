package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type HomeUI struct {
	Window fyne.Window
}

func NewHomeUI(window fyne.Window) *HomeUI {
	return &HomeUI{
		Window: window,
	}
}

func (hui *HomeUI) Home() fyne.CanvasObject {
	cardsContainer := container.NewGridWithColumns(4)

	return cardsContainer
}
