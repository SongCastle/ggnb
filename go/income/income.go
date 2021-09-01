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

func NewManager() (AbstractManager, error) {
	m := &Manager{}
	if err := m.Init(); err != nil {
		return nil, err
	}
	return m, nil
}

type AbstractManager interface {
	Init() error
	ToPayload(payload interface{}) (*bytes.Buffer, error)
	ToDummyPayload() (*bytes.Buffer, error)
}

type Manager struct {
	message message.Message
}

func (m *Manager) Init() error {
	switch os.Getenv(TypeEnv) {
	case GitHubType:
		m.message = &message.GitHubMessage{}
		return nil
	}
	return errors.New(fmt.Sprintf("Invalid %s", TypeEnv))
}

func (m *Manager) ToPayload(payload interface{}) (*bytes.Buffer, error) {
	if err := m.message.Init(payload); err != nil {
		return nil, err
	}
	return m.message.ToPayload()
}

func (m *Manager) ToDummyPayload() (*bytes.Buffer, error) {
	return m.message.ToDummyPayload()
}
