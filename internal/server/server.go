package server

import (
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
	sf.server = server.NewServer(dir, sf.LogCh)
	sf.server.Start()
}

func (sf *ServerFunctions) Receive(code string) error {
	return sf.server.Receive(code)
}

func (sf *ServerFunctions) Share(file string, callbacks server.ShareCallBacks) {
	sf.server.Share(file, callbacks)
}
