package handler

import (
	"github.com/SongCastle/ggnb/income"
	"github.com/SongCastle/ggnb/outcome"
	"github.com/aws/aws-lambda-go/lambda"
)

type lamdaHandler struct {}

func (_ *lamdaHandler) Start() {
	lambda.Start(
		func(payload interface{}) error {
			msg, err := income.ToPayload(payload)
			if err == nil {
				err = outcome.Send(msg)
			}
			return reportErrorIf(err)
		},
	)
}
