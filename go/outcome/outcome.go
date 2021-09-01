package outcome

import (
	"bytes"
	"fmt"

	"github.com/SongCastle/ggnb/income/builder"
)

func NewManager() AbstractManager {
	m := &Manager{}
	m.Init()
	return m
}

type AbstractManager interface {
	Init()
	Send(msg *bytes.Buffer) error
	ReportErrorIf(err error) error
}

type Manager struct {
	client client
}

func (m *Manager) Init() {
	m.client = &slackClient{}
	m.client.Init()
}

func (m *Manager) Send(msg *bytes.Buffer) error {
	if msg == nil {
		fmt.Println("Skipped")
		return nil
	}
	fmt.Printf("payload: %s\n", msg.String())

	if _, err := m.client.Post(msg); err != nil {
		return err
	}
	return nil
}

func (m *Manager) ReportErrorIf(err error) error {
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
		if msg, err := builder.BuildError(err); err != nil {
			fmt.Printf("Report Failed: %v\n", err)
		} else {
			if err := m.Send(msg); err != nil {
				fmt.Printf("Report Failed: %v\n", err)
			}
		}
	}
	return err
}
