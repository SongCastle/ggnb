package handler

import (
	"github.com/SongCastle/ggnb/income"
	"github.com/SongCastle/ggnb/outcome"
	"github.com/aws/aws-lambda-go/lambda"
)

type lambdaHandler struct {
	In income.AbstractManager
	Out outcome.AbstractManager
}

func (lh *lambdaHandler) Init(in income.AbstractManager, out outcome.AbstractManager) {
	lh.In = in
	lh.Out = out
}

func (lh *lambdaHandler) Start() {
	lambda.Start(
		func(payload interface{}) error {
			msg, err := lh.In.BuildMessage(payload)
			if err == nil {
				err = lh.Out.Send(msg)
			}
			return lh.Out.ReportErrorIf(err)
		},
	)
}
