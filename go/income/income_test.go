package income

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockedMessage struct{
  mock.Mock
}

func (m *mockedMessage) Init(payload interface{}) error {
	args := m.Called(payload)
  return args.Error(0)
}

func (m *mockedMessage) ToPayload() (*bytes.Buffer, error) {
	args := m.Called()
  return args[0].(*bytes.Buffer), args.Error(1)
}

func (m *mockedMessage) ToDummyPayload() (*bytes.Buffer, error) {
	args := m.Called()
  return args[0].(*bytes.Buffer), args.Error(1)
}

func TestNewManager(t *testing.T){
	assert := assert.New(t)
	beforeType := os.Getenv(TypeEnv)

	t.Run("without type", func(t *testing.T) {
		if err := os.Unsetenv(TypeEnv); err != nil {
			t.Fatal(err)
		}
		_, err := NewManager()
		assert.EqualError(err, fmt.Sprintf("Invalid %s", TypeEnv))
	})

	t.Run("GitHub", func(t *testing.T) {
		if err := os.Setenv(TypeEnv, GitHubType); err != nil {
			if err != nil {
				t.Fatal(err)
			}
		}
		m, err := NewManager()
		assert.Nil(err)
		assert.IsType(m, &Manager{})
	})

	t.Cleanup(func(){
		if err := os.Setenv(TypeEnv, beforeType); err != nil {
			if err != nil {
				t.Fatal(err)
			}
		}
	})
}

func TestManagerInit(t *testing.T) {
	assert := assert.New(t)
	beforeType := os.Getenv(TypeEnv)

	m := &Manager{}

	t.Run("without type", func(t *testing.T) {
		if err := os.Unsetenv(TypeEnv); err != nil {
			t.Fatal(err)
		}
		err := m.Init()
		assert.EqualError(err, fmt.Sprintf("Invalid %s", TypeEnv))
	})

	t.Run("GitHub", func(t *testing.T) {
		if err := os.Setenv(TypeEnv, GitHubType); err != nil {
			if err != nil {
				t.Fatal(err)
			}
		}
		err := m.Init()
		assert.Nil(err)
	})

	t.Cleanup(func(){
		if err := os.Setenv(TypeEnv, beforeType); err != nil {
			if err != nil {
				t.Fatal(err)
			}
		}
	})
}

func TestManagerToPayload(t *testing.T) {
	assert := assert.New(t)
	payload := map[string]string{"headers": "test", "body": "test"}

	t.Run("no errors", func(t *testing.T) {
		ebuf := bytes.NewBufferString(`{"body": "test"}`)
		mm := &mockedMessage{}
		mm.On("Init", payload).Return(nil)
		mm.On("ToPayload").Return(ebuf, nil)

		m := &Manager{message: mm}
		buf, err := m.ToPayload(payload)
		assert.Nil(err)
		assert.Equal(buf, ebuf)
	})

	t.Run("error", func(t *testing.T) {
		merr := "mocked"

		t.Run("Init", func(t *testing.T) {
			mm := &mockedMessage{}
			mm.On("Init", payload).Return(errors.New(merr))

			m := &Manager{message: mm}
			_, err := m.ToPayload(payload)
			assert.EqualError(err, merr)
		})

		t.Run("ToPayload", func(t *testing.T) {
			var ebuf *bytes.Buffer
			mm := &mockedMessage{}
			mm.On("Init", payload).Return(nil)
			mm.On("ToPayload").Return(ebuf, errors.New(merr))

			m := &Manager{message: mm}
			_, err := m.ToPayload(payload)
			assert.EqualError(err, merr)
		})
	})
}

func TestToDummyPayload(t *testing.T) {
	assert := assert.New(t)

	t.Run("no errors", func(t *testing.T) {
		ebuf := bytes.NewBufferString(`{"body": "test"}`)
		mm := &mockedMessage{}
		mm.On("ToDummyPayload").Return(ebuf, nil)

		m := &Manager{message: mm}
		buf, err := m.ToDummyPayload()
		assert.Nil(err)
		assert.Equal(buf, ebuf)
	})

	t.Run("error", func(t *testing.T) {
		merr := "mocked"
		var ebuf *bytes.Buffer
		mm := &mockedMessage{}
		mm.On("ToDummyPayload").Return(ebuf, errors.New(merr))

		m := &Manager{message: mm}
		_, err := m.ToDummyPayload()
		assert.EqualError(err, merr)
	})
}
