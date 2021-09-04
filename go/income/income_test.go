package income

import (
	"bytes"
	"errors"
	"testing"

	"github.com/SongCastle/ggnb/income/message"
	"github.com/stretchr/testify/assert"
)

func TestNewManager(t *testing.T){
	t.Parallel()

	m := NewManager(&message.MockedMessage{})
	_, ok := m.(AbstractManager)
	assert.True(t, ok)
}

func TestManagerInit(t *testing.T) {
	t.Parallel()

	m := &Manager{}
	assert.NotPanics(
		t,
		func() { m.Init(&message.MockedMessage{}) },
	)
}

func TestManagerBuildMessage(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	payload := map[string]string{"headers": "test", "body": "test"}

	t.Run("no errors", func(t *testing.T) {
		ebuf := bytes.NewBufferString(`{"body": "test"}`)
		mm := &message.MockedMessage{}
		mm.On("Init", payload).Return(nil)
		mm.On("ToPayload").Return(ebuf, nil)

		m := &Manager{message: mm}
		buf, err := m.BuildMessage(payload)
		assert.Nil(err)
		assert.Equal(buf, ebuf)
	})

	t.Run("error", func(t *testing.T) {
		eemsg := "mocked"
		eerr := errors.New(eemsg)

		t.Run("Init", func(t *testing.T) {
			mm := &message.MockedMessage{}
			mm.On("Init", payload).Return(eerr)

			m := &Manager{message: mm}
			_, err := m.BuildMessage(payload)
			assert.EqualError(err, eemsg)
		})

		t.Run("ToPayload", func(t *testing.T) {
			var ebuf *bytes.Buffer
			mm := &message.MockedMessage{}
			mm.On("Init", payload).Return(nil)
			mm.On("ToPayload").Return(ebuf, eerr)

			m := &Manager{message: mm}
			_, err := m.BuildMessage(payload)
			assert.EqualError(err, eemsg)
		})
	})
}

func TestBuildDummyMessage(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	t.Run("no errors", func(t *testing.T) {
		ebuf := bytes.NewBufferString(`{"body": "test"}`)
		mm := &message.MockedMessage{}
		mm.On("ToDummyPayload").Return(ebuf, nil)

		m := &Manager{message: mm}
		buf, err := m.BuildDummyMessage()
		assert.Nil(err)
		assert.Equal(buf, ebuf)
	})

	t.Run("error", func(t *testing.T) {
		var ebuf *bytes.Buffer
		eemsg := "mocked"
		mm := &message.MockedMessage{}
		mm.On("ToDummyPayload").Return(ebuf, errors.New(eemsg))

		m := &Manager{message: mm}
		_, err := m.BuildDummyMessage()
		assert.EqualError(err, eemsg)
	})
}
