package outcome

import (
	"bytes"
	"errors"
	"testing"

	"github.com/SongCastle/ggnb/outcome/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewManager(t *testing.T) {
	t.Parallel()
	m := NewManager(&client.MockedClient{})
	_, ok := m.(AbstractManager)
	assert.True(t, ok)
}

func TestManagerInit(t *testing.T) {
	t.Parallel()

	c := &client.MockedClient{}
	m := &Manager{}
	m.Init(c)
	assert.IsType(t, m.client, c)
}

func TestManagerSend(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	msg := bytes.NewBufferString(`{"body": "test"}`)

	t.Run("no errors", func(t *testing.T) {
		c := &client.MockedClient{}
		c.On("Post", msg).Return([]byte("ok"), nil)

		m := &Manager{client: c}
		err := m.Send(msg)
		assert.Nil(err)
	})

	t.Run("error", func(t *testing.T) {
		var b []byte
		eemsg := "mocked"

		c := &client.MockedClient{}
		c.On("Post", msg).Return(b, errors.New(eemsg))

		m := &Manager{client: c}
		err := m.Send(msg)
		assert.EqualError(err, eemsg)
	})
}

func TestManagerReportErrorIf(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	t.Run("no errors", func(t *testing.T) {
		c := &client.MockedClient{}
		c.On("Post", mock.AnythingOfType("*bytes.Buffer")).Return([]byte("ok"), nil)

		m := &Manager{client: c}
		err := m.ReportErrorIf(nil)
		assert.Nil(err)
	})

	t.Run("error", func(t *testing.T) {
		var b []byte
		eemsg, eemsg2 := "mocked", "mocked2"

		c := &client.MockedClient{}
		c.On("Post", mock.AnythingOfType("*bytes.Buffer")).Return(b, errors.New(eemsg2))

		m := &Manager{client: c}
		err := m.ReportErrorIf(errors.New(eemsg))
		assert.EqualError(err, eemsg)
	})
}
