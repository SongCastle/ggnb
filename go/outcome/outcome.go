package outcome

import (
	"bytes"
	"fmt"

	"github.com/SongCastle/ggnb/income/builder"
	"github.com/SongCastle/ggnb/outcome/client"
)

func NewManager(client client.AbstractClient) AbstractManager {
	m := &Manager{}
	m.Init(client)
	return m
}

type AbstractManager interface {
	Init(client.AbstractClient)
	Send(*bytes.Buffer) error
	ReportErrorIf(error) error
}

type Manager struct {
	client client.AbstractClient
}

func (m *Manager) Init(client client.AbstractClient) {
	m.client = client
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
