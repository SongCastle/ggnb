package income

import (
	"bytes"

	"github.com/SongCastle/ggnb/income/message"
	"github.com/stretchr/testify/mock"
)

type MockedIncomeManager struct {
	mock.Mock
}

func (im *MockedIncomeManager) Init(message message.AbstractMessage) {
}

func (im *MockedIncomeManager) BuildMessage(headers, body interface{}) (*bytes.Buffer, error) {
	args := im.Called(headers, body)
	return args[0].(*bytes.Buffer), args.Error(1)
}

func (im *MockedIncomeManager) BuildDummyMessage() (*bytes.Buffer, error) {
	args := im.Called()
	return args[0].(*bytes.Buffer), args.Error(1)
}
