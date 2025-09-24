package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/Owbird/SNetT-Engine/pkg/config"
	"github.com/Owbird/SNetT-Engine/pkg/models"
	engineServer "github.com/Owbird/SNetT-Engine/pkg/server"
	"github.com/owbird/snett/internal/server"
	"github.com/skratchdot/open-golang/open"
)

type LogStatus = string

const (
	LogSuccess LogStatus = "success"
	LogError   LogStatus = "error"
)

type ServerUI struct {
	Window    fyne.Window
	Functions *server.ServerFunctions
}

func NewServerUI(window fyne.Window) *ServerUI {
	return &ServerUI{
		Window:    window,
		Functions: server.NewServerFunctions(),
	}
}

func (sui *ServerUI) ReceiveFile() {
	codeInput := widget.NewEntry()
	codeInput.SetPlaceHolder("2-code-here")

	codeForm := widget.NewFormItem("Code", codeInput)

	formItems := []*widget.FormItem{
		codeForm,
	}

	callback := func(create bool) {
		if create {
			if err := sui.Functions.Receive(codeInput.Text); err != nil {
				dialog.NewError(err, sui.Window)
			}
		}
	}

	dialog.NewForm("Enter code",
		"Receive",
		"Cancel",
		formItems,
		callback,
		sui.Window,
	).Show()
}

func (sui *ServerUI) ShareFile() {
	dialog.ShowFileOpen(func(uc fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.NewError(err, sui.Window)
			return
		}

		file := uc.URI().Path()

		var codeReceivedDialog dialog.Dialog

		sui.Functions.Share(file, engineServer.ShareCallBacks{
			OnFileSent: func() {
				codeReceivedDialog.Hide()
				dialog.NewInformation("File sent", "File sent successfully", sui.Window).Show()
			},
			OnSendErr: func(err error) {
				dialog.NewError(err, sui.Window).Show()
			},
			OnProgressChange: func(progress models.FileShareProgress) {},
			OnCodeReceive: func(code string) {
				codeReceivedDialog = dialog.NewInformation("Code received", code, sui.Window)

				codeReceivedDialog.Show()
			},
		})
	}, sui.Window)
}

func (sui *ServerUI) ChooseHostDir() {
	dialog.ShowFolderOpen(func(lu fyne.ListableURI, err error) {
		if err != nil {
			dialog.NewError(err, sui.Window)
			return
		}

		go sui.Functions.Host(lu.Path())

		logWindow := fyne.CurrentApp().NewWindow("Server Logs")
		logWindow.Resize(fyne.NewSize(500, 500))

		logsContainer := container.NewVBox()

		logWindow.SetContent(container.NewVScroll(logsContainer))

		logWindow.Show()

		showLog := func(status LogStatus, text string) {
			if status == LogError {
				logsContainer.Add(
					widget.NewRichText(
						&widget.TextSegment{
							Text: text,
							Style: widget.RichTextStyle{
								ColorName: "red",
							},
						},
					),
				)
			} else {
				logsContainer.Add(
					widget.NewRichText(&widget.TextSegment{Text: text}))
			}
		}

		go func() {
			for l := range sui.Functions.LogCh {
				switch l.Type {
				case models.API_LOG:
					if l.Error != nil {
						showLog(LogError, fmt.Sprintf("[!] API Log [error]: %v", l.Error.Error()))
					} else {
						showLog(LogSuccess, fmt.Sprintf("[+] API Log: %v", l.Message))
					}

				case models.SERVE_WEB_UI_NETWORK:
					showLog(LogSuccess, fmt.Sprintf("[+] Network Web Running: %v", l.Message))

				case models.SERVE_WEB_UI_REMOTE:
					showLog(LogSuccess, fmt.Sprintf("[+] Remote Web Running: %v", l.Message))

					open.Run(l.Message)
					open.Run("https://loca.lt/mytunnelpassword")

				default:
					if l.Error != nil {
						showLog(LogError, fmt.Sprintf("[!] Error: %v", l.Error.Error()))
					} else {
						showLog(LogError, fmt.Sprintf("[+] Log: %v", l.Message))
					}

				}
			}
			logsContainer.Refresh()
		}()
	}, sui.Window)
}

func (sui *ServerUI) ServerSettings() {
	settingsWindow := fyne.CurrentApp().NewWindow("Server settings")
	settingsWindow.Resize(fyne.NewSize(500, 500))

	appConfig := config.NewAppConfig()

	serverConfig := appConfig.GetSeverConfig()

	serverNameInput := widget.NewEntry()
	serverNameInput.SetPlaceHolder("Enter name")

	allowUploadsChecker := widget.NewCheck("Allow uploads", func(value bool) {
		serverConfig.SetAllowUploads(value)
	})

	serverNameInput.Text = serverConfig.GetName()
	allowUploadsChecker.Checked = serverConfig.GetAllowUploads()

	saveBtn := widget.NewButton("Save", func() {
		serverConfig.SetName(serverNameInput.Text)

		appConfig.Save()

		settingsWindow.Close()
	})

	layoutContainer := container.NewVBox(serverNameInput, allowUploadsChecker, saveBtn)

	settingsWindow.SetContent(layoutContainer)

	settingsWindow.Show()
}
