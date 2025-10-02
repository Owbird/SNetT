package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/Owbird/SNetT-Engine/pkg/models"
	engineWormhole "github.com/Owbird/SNetT-Engine/pkg/wormhole"
	"github.com/owbird/snett/internal/wormhole"
)

type WormholeUI struct {
	Window    fyne.Window
	Functions *wormhole.WormholeFunctions
}

func NewWormholeUI(window fyne.Window) *WormholeUI {
	return &WormholeUI{
		Window:    window,
		Functions: wormhole.NewWormholeFunctions(),
	}
}

func (wui *WormholeUI) ReceiveFile() {
	codeInput := widget.NewEntry()
	codeInput.SetPlaceHolder("2-code-here")

	codeForm := widget.NewFormItem("Code", codeInput)

	formItems := []*widget.FormItem{
		codeForm,
	}

	callback := func(create bool) {
		if create {
			if err := wui.Functions.Receive(codeInput.Text); err != nil {
				dialog.NewError(err, wui.Window)
			}
		}
	}

	dialog.NewForm("Enter code",
		"Receive",
		"Cancel",
		formItems,
		callback,
		wui.Window,
	).Show()
}

func (wui *WormholeUI) ShareFile() {
	dialog.ShowFileOpen(func(uc fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.NewError(err, wui.Window)
			return
		}

		file := uc.URI().Path()

		var codeReceivedDialog dialog.Dialog

		wui.Functions.Share(file, engineWormhole.ShareCallBacks{
			OnFileSent: func() {
				codeReceivedDialog.Hide()
				dialog.NewInformation("File sent", "File sent successfully", wui.Window).Show()
			},
			OnSendErr: func(err error) {
				dialog.NewError(err, wui.Window).Show()
			},
			OnProgressChange: func(progress models.FileShareProgress) {},
			OnCodeReceive: func(code string) {
				codeReceivedDialog = dialog.NewInformation("Code received", code, wui.Window)

				codeReceivedDialog.Show()
			},
		})
	}, wui.Window)
}
