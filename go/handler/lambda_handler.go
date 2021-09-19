package handler

import (
	"github.com/SongCastle/ggnb/income"
	"github.com/SongCastle/ggnb/outcome"

	"github.com/aws/aws-lambda-go/events"
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
		func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
			msg, err := lh.In.BuildMessage(request.Headers, &request.Body)
			if err == nil {
				err = lh.Out.Send(msg)
			}

			body, statusCode := "ok", 200
			if err != nil {
				body, statusCode = err.Error(), 400
				lh.Out.ReportErrorIf(err)
			}
			return events.APIGatewayProxyResponse{Body: body, StatusCode: statusCode}, nil
		},
	)
}
