package income

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/SongCastle/ggnb/income/message"
)

const (
	TypeEnv = "INCOME_TYPE"
	GitHubType = "github"
)

func newMessage() (message.Message, error) {
	switch os.Getenv(TypeEnv) {
	case GitHubType:
		return &message.GitHubMessage{}, nil
	}
	return nil, errors.New(fmt.Sprintf("Invalid %s", TypeEnv))
}

func ToPayload(payload interface{}) (*bytes.Buffer, error) {
	msg, err := newMessage()
	if err != nil {
		return nil, err
	}
	if err := msg.Init(payload); err != nil {
		return nil, err
	}
	return msg.ToPayload()
}

func ToDummyPayload() (*bytes.Buffer, error) {
	msg, err := newMessage()
	if err != nil {
		return nil, err
	}
	return msg.ToDummyPayload()
}
