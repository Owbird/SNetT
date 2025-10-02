package wormhole

import (
	"github.com/Owbird/SNetT-Engine/pkg/wormhole"
)

type WormholeFunctions struct {
	wormhole *wormhole.Wormhole
}

func NewWormholeFunctions() *WormholeFunctions {

	return &WormholeFunctions{}
}

func (wf *WormholeFunctions) Receive(code string) error {
	return wf.wormhole.Receive(code)

	return nil
}

func (wf *WormholeFunctions) Share(file string, callbacks wormhole.ShareCallBacks) {
	wf.wormhole.Share(file, callbacks)
}
