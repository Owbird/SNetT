package ui

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/Owbird/SNetT-Engine/pkg/config"
	"github.com/Owbird/SNetT-Engine/pkg/models"

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

func (wui *ServerUI) ChooseHostDir() {
	dialog.ShowFolderOpen(func(lu fyne.ListableURI, err error) {
		if err != nil {
			dialog.NewError(err, wui.Window)
			return
		}

		if lu == nil {
			return
		}

		go wui.Functions.Host(lu.Path())

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
			for l := range wui.Functions.LogCh {
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
	}, wui.Window)
}

func (wui *ServerUI) ServerSettings() {
	settingsWindow := fyne.CurrentApp().NewWindow("Server settings")
	settingsWindow.Resize(fyne.NewSize(500, 500))

	appConfig := config.NewAppConfig()

	serverConfig := appConfig.GetSeverConfig()

	serverNameInput := widget.NewEntry()
	serverNameInput.SetPlaceHolder("Enter name")

	serverPort := widget.NewEntry()
	serverPort.SetPlaceHolder("Enter port")

	allowUploadsChecker := widget.NewCheck("Allow uploads", func(value bool) {
		serverConfig.SetAllowUploads(value)
	})

	allowOnlineChecker := widget.NewCheck("Allow online", func(value bool) {
		serverConfig.SetAllowOnline(value)
	})

	serverNameInput.Text = serverConfig.GetName()
	serverNameInput.OnChanged = func(s string) {
		serverConfig.SetName(s)
	}

	serverPort.Text = strconv.Itoa(serverConfig.GetPort())
	serverPort.OnChanged = func(s string) {
		port, _ := strconv.Atoi(s)
		serverConfig.SetPort(port)
	}

	allowUploadsChecker.Checked = serverConfig.GetAllowUploads()
	allowOnlineChecker.Checked = serverConfig.GetAllowOnline()

	saveBtn := widget.NewButton("Save", func() {
		appConfig.Save()

		settingsWindow.Close()
	})

	layoutContainer := container.NewVBox(serverNameInput, serverPort, allowUploadsChecker, allowOnlineChecker, saveBtn)

	settingsWindow.SetContent(layoutContainer)

	settingsWindow.Show()
}
