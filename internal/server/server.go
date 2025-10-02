package server

import (
	"github.com/Owbird/SNetT-Engine/pkg/config"
	"github.com/Owbird/SNetT-Engine/pkg/models"
	"github.com/Owbird/SNetT-Engine/pkg/server"
)

type ServerFunctions struct {
	server *server.Server
	LogCh  chan models.ServerLog
}

func NewServerFunctions() *ServerFunctions {
	logCh := make(chan models.ServerLog)

	return &ServerFunctions{
		LogCh: logCh,
	}
}

func (sf *ServerFunctions) Host(dir string) {
	appConfig := config.NewAppConfig()

	sf.server = server.NewServer(dir, sf.LogCh)
	sf.server.Start(*appConfig)
}

func (sf *ServerFunctions) Receive(code string) error {
	// return sf.server.Receive(code)

	return nil
}

func (sf *ServerFunctions) Share(file string, callbacks any) {
	// sf.server.Share(file, callbacks)
}
