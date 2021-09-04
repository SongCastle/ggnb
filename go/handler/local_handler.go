package handler

import (
	"github.com/SongCastle/ggnb/income"
	"github.com/SongCastle/ggnb/outcome"
)

type localHandler struct {
	In income.AbstractManager
	Out outcome.AbstractManager
}

func (lh *localHandler) Init(in income.AbstractManager, out outcome.AbstractManager) {
	lh.In = in
	lh.Out = out
}

func (lh *localHandler) Start() {
	msg, err := lh.In.BuildDummyMessage()
	if err == nil {
		err = lh.Out.Send(msg)
	}
	lh.Out.ReportErrorIf(err)
}
