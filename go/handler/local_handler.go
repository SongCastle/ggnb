package handler

import (
	"github.com/SongCastle/ggnb/income"
	"github.com/SongCastle/ggnb/outcome"
)

type localHandler struct {}

func (_ *localHandler) Start() {
	msg, err := income.ToDummyPayload()
	if err == nil {
		err = outcome.Send(msg)
	}
	reportErrorIf(err)
}
