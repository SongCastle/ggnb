package outcome

import (
	"bytes"

	"github.com/SongCastle/ggnb/outcome/client"
	"github.com/stretchr/testify/mock"
)

type MockedOutcomeManager struct {
	mock.Mock
}

func (om *MockedOutcomeManager) Init(_ client.AbstractClient) {
}

func (om *MockedOutcomeManager) Send(msg *bytes.Buffer) error {
	args := om.Called(msg)
	return args.Error(0)
}

func (om *MockedOutcomeManager) ReportErrorIf(err error) error {
	args := om.Called(err)
	return args.Error(0)
}
