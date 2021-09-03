package handler

import (
	"os"

	"github.com/SongCastle/ggnb/income"
	"github.com/SongCastle/ggnb/income/message"
	"github.com/SongCastle/ggnb/outcome"
	"github.com/SongCastle/ggnb/outcome/client"
)

const DEBUG = "DEBUG"

type abstractHandler interface {
	Init(income.AbstractManager, outcome.AbstractManager)
	Start()
}

func New(am message.AbstractMessage, ac client.AbstractClient) abstractHandler {
	var h abstractHandler
	if onLambda() {
		h = &lambdaHandler{}
	} else {
		h = &localHandler{}
	}
	// Create Manager
	in := income.NewManager(am)
	out := outcome.NewManager(ac)
	h.Init(in, out)
	return h
}

func onLambda() bool {
	return os.Getenv(DEBUG) == ""
}
