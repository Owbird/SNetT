package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"net/url"
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
	name := widget.NewLabel("SNetT v2.0")
	name.Alignment = fyne.TextAlignCenter
	name.TextStyle = fyne.TextStyle{Bold: true}

	intro := widget.NewLabel("A secure, cross-platform desktop application for managing files over a network.")
	intro.Wrapping = fyne.TextWrapWord
	intro.Alignment = fyne.TextAlignCenter

	githubURL, _ := url.Parse("https://github.com/Owbird/SNetT")
	githubLink := widget.NewHyperlink("View on Github", githubURL)
	githubLink.Alignment = fyne.TextAlignCenter

	content := container.New(
		layout.NewVBoxLayout(),
		layout.NewSpacer(),
		name,
		intro,
		githubLink,
		layout.NewSpacer(),
	)

	return content
}
