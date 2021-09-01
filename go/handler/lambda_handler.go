package handler

import (
	"fmt"

	"github.com/SongCastle/ggnb/income"
	"github.com/SongCastle/ggnb/outcome"
	"github.com/aws/aws-lambda-go/lambda"
)

type lambdaHandler struct {
	in income.AbstractManager
	out outcome.AbstractManager
}

func (lh *lambdaHandler) Init() error {
	var err error
	lh.in, err = income.NewManager()
	if err != nil {
		fmt.Printf("Init failed: %v\n", err)
		return err
	}
	lh.out = outcome.NewManager()
	return lh.out.ReportErrorIf(err)
}

func (lh *lambdaHandler) Start() {
	if lh.in != nil && lh.out != nil {
		lambda.Start(
			func(payload interface{}) error {
				msg, err := lh.in.ToPayload(payload)
				if err == nil {
					err = lh.out.Send(msg)
				}
				return lh.out.ReportErrorIf(err)
			},
		)
	}
}
