package message

import (
	"bytes"
	"fmt"
	"errors"
	"os"
)

const (
	TypeEnv = "INCOME_TYPE"
	GitHubType = "github"
)

type AbstractMessage interface {
	Init(payload interface{}) error
	ToPayload() (*bytes.Buffer, error)
	ToDummyPayload() (*bytes.Buffer, error)
}

func NewMessage() (AbstractMessage, error) {
	// Only GitHub
	switch os.Getenv(TypeEnv) {
	case GitHubType:
		return &GitHubMessage{}, nil
	}
	return nil, errors.New(fmt.Sprintf("Invalid %s", TypeEnv))
}
