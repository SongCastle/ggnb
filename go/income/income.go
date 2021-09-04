package income

import (
	"bytes"

	"github.com/SongCastle/ggnb/income/message"
)

func NewManager(msg message.AbstractMessage) AbstractManager {
	m := &Manager{}
	m.Init(msg)
	return m
}

type AbstractManager interface {
	Init(message message.AbstractMessage)
	BuildMessage(payload interface{}) (*bytes.Buffer, error)
	BuildDummyMessage() (*bytes.Buffer, error)
}

type Manager struct {
	message message.AbstractMessage
}

func (m *Manager) Init(message message.AbstractMessage) {
	m.message = message
}

func (m *Manager) BuildMessage(payload interface{}) (*bytes.Buffer, error) {
	if err := m.message.Init(payload); err != nil {
		return nil, err
	}
	return m.message.ToPayload()
}

func (m *Manager) BuildDummyMessage() (*bytes.Buffer, error) {
	return m.message.ToDummyPayload()
}
