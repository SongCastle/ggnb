package client

import (
	"bytes"

	"github.com/stretchr/testify/mock"
)

type MockedClient struct {
  mock.Mock
}

func (m *MockedClient) Init() error {
	args := m.Called()
  return args.Error(0)
}

func (m *MockedClient) Post(msg *bytes.Buffer) ([]byte, error) {
	args := m.Called(msg)
  return args[0].([]byte), args.Error(1)
}
