package outcome

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockedClient struct {
  mock.Mock
}

func (_ *mockedClient) Init() {
}

func (m *mockedClient) Post(msg *bytes.Buffer) ([]byte, error) {
	args := m.Called(msg)
  return args[0].([]byte), args.Error(1)
}

func TestNewManager(t *testing.T) {
	t.Parallel()
	m := NewManager()
	assert.IsType(t, m, &Manager{})
}

func TestManagerInit(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	m := &Manager{}
	m.Init()
	assert.IsType(m.client, &slackClient{})
	assert.NotNil(m.client)
}

func TestManagerSend(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	msg := bytes.NewBufferString(`{"body": "test"}`)

	t.Run("no errors", func(t *testing.T) {
		mc := &mockedClient{}
		mc.On("Post", msg).Return([]byte("ok"), nil)

		m := &Manager{client: mc}
		err := m.Send(msg)
		assert.Nil(err)
	})

	t.Run("error", func(t *testing.T) {
		merr := "mocked"
		var b []byte
		mc := &mockedClient{}
		mc.On("Post", msg).Return(b, errors.New(merr))

		m := &Manager{client: mc}
		err := m.Send(msg)
		assert.EqualError(err, merr)
	})
}

func TestManagerReportErrorIf(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	t.Run("no errors", func(t *testing.T) {
		mc := &mockedClient{}
		mc.On("Post", mock.AnythingOfType("*bytes.Buffer")).Return([]byte("ok"), nil)

		m := &Manager{client: mc}
		err := m.ReportErrorIf(nil)
		assert.Nil(err)
	})

	t.Run("error", func(t *testing.T) {
		merr, merr2 := "mocked", "mocked2"
		var b []byte
		mc := &mockedClient{}
		mc.On("Post", mock.AnythingOfType("*bytes.Buffer")).Return(b, errors.New(merr2))

		m := &Manager{client: mc}
		err := m.ReportErrorIf(errors.New(merr))
		assert.EqualError(err, merr)
	})
}
