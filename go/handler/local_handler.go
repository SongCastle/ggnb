package handler

import (
	"fmt"

	"github.com/SongCastle/ggnb/income"
	"github.com/SongCastle/ggnb/outcome"
)

type localHandler struct {
	in income.AbstractManager
	out outcome.AbstractManager
}

func (lh *localHandler) Init() error {
	var err error
	lh.in, err = income.NewManager()
	if err != nil {
		fmt.Printf("Init failed: %v\n", err)
		return err
	}
	lh.out = outcome.NewManager()
	return nil
}

func (lh *localHandler) Start() {
	if lh.in != nil && lh.out != nil {
		msg, err := lh.in.ToDummyPayload()
		if err == nil {
			err = lh.out.Send(msg)
		}
		lh.out.ReportErrorIf(err)
	}
}
