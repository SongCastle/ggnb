package main

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/SongCastle/ggnb/income"
	"github.com/SongCastle/ggnb/outcome"
	"github.com/aws/aws-lambda-go/lambda"
)

const DEBUG = "DEBUG"

func receive(payload interface{}) (*bytes.Buffer, error) {
	_payload, err := income.ToPayload(payload)
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
		return nil, err
	}
	return _payload, nil
}

func receiveDummy() (*bytes.Buffer, error) {
	_payload, err := income.ToDummyPayload()
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
		return nil, err
	}
	return _payload, nil
}

func send(msg *bytes.Buffer) error {
	if msg == nil {
		fmt.Println("Skipped\n")
		return nil
	}
	fmt.Printf("payload: %s\n", msg.String())

	out := outcome.New()
	err := out.Post(msg)
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
		return err
	}
	return nil
}

func HandleRequest(ctx context.Context, payload interface{}) error {
	msg, err := receive(payload)
	if err != nil {
		return err
	}
	return send(msg)
}

func HandleDummyRequest() error {
	msg, err := receiveDummy()
	if err != nil {
		return err
	}
	return send(msg)
}

func onLambda() bool {
	return os.Getenv(DEBUG) == ""
}

func main() {
	if onLambda() {
		lambda.Start(HandleRequest)
		return
	}

	HandleDummyRequest()
}
