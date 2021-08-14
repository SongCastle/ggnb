package income

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/SongCastle/ggnb/message"
)

const (
	GitHubType = "github"
	TypeEnv = "INCOME_TYPE"
)

func New() (message.Message, error) {
	switch os.Getenv(TypeEnv) {
	case GitHubType:
		return &message.GitHubMessage{}, nil
	}
	return nil, errors.New(fmt.Sprintf("Invalid %s", TypeEnv))
}

func ToPayload(payload interface{}) (*bytes.Buffer, error) {
	msg, err := New()
	if err != nil {
		return nil, err
	}
	if err := msg.Init(payload); err != nil {
		return nil, err
	}
	if !msg.NeedToDeliver() {
		return nil, nil
	}
	return msg.ToPayload(), nil
}

func ToDummyPayload() (*bytes.Buffer, error) {
	msg, err := New()
	if err != nil {
		return nil, err
	}
	return msg.ToDummyPayload(), nil
}
