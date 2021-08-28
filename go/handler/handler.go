package handler

import (
	"fmt"
	"os"

	"github.com/SongCastle/ggnb/income/builder"
	"github.com/SongCastle/ggnb/outcome"
)

const DEBUG = "DEBUG"

type Handler interface {
	Start()
}

func New() Handler {
	if onLambda() {
		return &lamdaHandler{}
	}
	return &localHandler{}
}

func onLambda() bool {
	return os.Getenv(DEBUG) == ""
}

func reportErrorIf(err error) error {
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
		if msg, err := builder.BuildError(err); err != nil {
			fmt.Printf("Report Failed: %v\n", err)
		} else {
			if err := outcome.Send(msg); err != nil {
				fmt.Printf("Report Failed: %v\n", err)
			}
		}
	}
	return err
}
