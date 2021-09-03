package message

import (
	"bytes"

	"github.com/stretchr/testify/mock"
)

type MockedMessage struct{
  mock.Mock
}

func (m *MockedMessage) Init(payload interface{}) error {
	args := m.Called(payload)
  return args.Error(0)
}

func (m *MockedMessage) ToPayload() (*bytes.Buffer, error) {
	args := m.Called()
  return args[0].(*bytes.Buffer), args.Error(1)
}

func (m *MockedMessage) ToDummyPayload() (*bytes.Buffer, error) {
	args := m.Called()
  return args[0].(*bytes.Buffer), args.Error(1)
}
