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

	"github.com/Owbird/SNetT-Engine/pkg/server"
	"github.com/skratchdot/open-golang/open"
)

type LogStatus = string

const (
	LogSuccess LogStatus = "success"
	LogError   LogStatus = "error"
)

type ServerUI struct {
	Window fyne.Window
	Server *server.Server
	LogCh  chan models.ServerLog
}

func NewServerUI(window fyne.Window) *ServerUI {
	logCh := make(chan models.ServerLog)

	return &ServerUI{
		Window: window,
		LogCh:  logCh,
		Server: server.NewServer("", logCh),
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

		appConfig := config.NewAppConfig()

		wui.Server.Dir = lu.Path()

		go wui.Server.Start(*appConfig)

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
			for l := range wui.LogCh {
				switch l.Type {
				case models.SERVER_ERROR:
					showLog(LogError, fmt.Sprintf("[!] API Log [error]: %v", l.Value))

				case models.API_LOG:
					showLog(LogSuccess, fmt.Sprintf("[+] %v", l.Value))

				case models.SERVE_UI_LOCAL:
					showLog(LogSuccess, fmt.Sprintf("[+] Network Web Running: %v", l.Value))

				case models.SERVE_UI_REMOTE:
					showLog(LogSuccess, fmt.Sprintf("[+] Remote Web Running: %v", l.Value))

					open.Run(l.Value)
					open.Run("https://loca.lt/mytunnelpassword")

				default:
					showLog(LogError, fmt.Sprintf("[+] Log: %v", l.Value))

				}
			}
			logsContainer.Refresh()
		}()
	}, wui.Window)
}

func (wui *ServerUI) Discover() {
	logWindow := fyne.CurrentApp().NewWindow("Available servers")
	logWindow.Resize(fyne.NewSize(500, 500))

	logsContainer := container.NewVBox()

	servers := make(chan models.SNetTServer)

	logWindow.SetContent(container.NewVScroll(logsContainer))

	logWindow.Show()

	go wui.Server.List(servers)

	for s := range servers {
		logsContainer.Add(
			widget.NewRichText(
				&widget.TextSegment{
					Text: fmt.Sprintf("[%v] %v:%v", s.Name, s.IP, s.Port),
				},
			),
		)
	}
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
		serverConfig.AllowUploads = value
	})

	allowOnlineChecker := widget.NewCheck("Allow online", func(value bool) {
		serverConfig.AllowOnline = value
	})

	serverNameInput.Text = serverConfig.Name
	serverNameInput.OnChanged = func(s string) {
		serverConfig.Name = s
	}

	serverPort.Text = strconv.Itoa(serverConfig.Port)
	serverPort.OnChanged = func(s string) {
		port, _ := strconv.Atoi(s)
		serverConfig.Port = port
	}

	allowUploadsChecker.Checked = serverConfig.AllowUploads
	allowOnlineChecker.Checked = serverConfig.AllowOnline

	saveBtn := widget.NewButton("Save", func() {
		appConfig.Save()

		settingsWindow.Close()
	})

	layoutContainer := container.NewVBox(serverNameInput, serverPort, allowUploadsChecker, allowOnlineChecker, saveBtn)

	settingsWindow.SetContent(layoutContainer)

	settingsWindow.Show()
}
