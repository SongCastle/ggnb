package handler

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/SongCastle/ggnb/income"
	"github.com/SongCastle/ggnb/outcome"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockedIncomeManager struct {
	mock.Mock
}

func (im *mockedIncomeManager) Init() error {
	args := im.Called()
	return args.Error(0)
}

func (im *mockedIncomeManager) ToPayload(payload interface{}) (*bytes.Buffer, error) {
	args := im.Called(payload)
	return args[0].(*bytes.Buffer), args.Error(1)
}

func (im *mockedIncomeManager) ToDummyPayload() (*bytes.Buffer, error) {
	args := im.Called()
	return args[0].(*bytes.Buffer), args.Error(1)
}

type mockedOutcomeManager struct {
	mock.Mock
}

func (om *mockedOutcomeManager) Init() {
}

func (om *mockedOutcomeManager) Send(msg *bytes.Buffer) error {
	args := om.Called(msg)
	return args.Error(0)
}

func (om *mockedOutcomeManager) ReportErrorIf(err error) error {
	args := om.Called(err)
	return args.Error(0)
}

func TestHandlerNew(t *testing.T) {
	assert := assert.New(t)
	beforeDebug := os.Getenv(DEBUG)

	t.Run("debug", func(t *testing.T) {
		if err := os.Setenv(DEBUG, "1"); err != nil {
			t.Fatal(err)
		}
		h := New()
		assert.IsType(h, &localHandler{})
	})

	t.Run("lambda", func(t *testing.T) {
		if err := os.Unsetenv(DEBUG); err != nil {
			t.Fatal(err)
		}
		h := New()
		assert.IsType(h, &lambdaHandler{})
	})

	t.Cleanup(func(){
		if err := os.Setenv(DEBUG, beforeDebug); err != nil {
			if err != nil {
				t.Fatal(err)
			}
		}
	})
}

func TestLocalHandlerInit(t *testing.T) {
	assert := assert.New(t)
	beforeType := os.Getenv(income.TypeEnv)

	t.Run("no errors", func(t *testing.T) {
		if err := os.Setenv(income.TypeEnv, income.GitHubType); err != nil {
			t.Fatal(err)
		}

		h := &localHandler{}
		err := h.Init()
		assert.Nil(err)
		assert.IsType(h.in, &income.Manager{})
		assert.IsType(h.out, &outcome.Manager{})
	})

	t.Run("error", func(t *testing.T) {
		if err := os.Unsetenv(income.TypeEnv); err != nil {
			t.Fatal(err)
		}

		h := &localHandler{}
		err := h.Init()
		assert.EqualError(err, fmt.Sprintf("Invalid %s", income.TypeEnv))
		assert.Nil(h.in)
		assert.Nil(h.out)
	})

	t.Cleanup(func(){
		if err := os.Setenv(income.TypeEnv, beforeType); err != nil {
			if err != nil {
				t.Fatal(err)
			}
		}
	})
}

func TestLocalHandlerStart(t *testing.T) {
	t.Run("no errors", func(t *testing.T) {
		msg := bytes.NewBufferString(`{"body": "test"}`)
		min := &mockedIncomeManager{}
		min.On("ToDummyPayload").Return(msg, nil)

		mout := &mockedOutcomeManager{}
		mout.On("Send", msg).Return(nil)
		mout.On("ReportErrorIf", nil).Return(nil)

		h := &localHandler{in: min, out: mout}
		h.Start()
	})

	t.Run("error", func(t *testing.T) {
		memsg := "mocked"
		merr := errors.New(memsg)

		var b *bytes.Buffer
		min := &mockedIncomeManager{}
		min.On("ToDummyPayload").Return(b, merr)

		mout := &mockedOutcomeManager{}
		mout.On("ReportErrorIf", merr).Return(merr)

		h := &localHandler{in: min, out: mout}
		h.Start()
	})

	t.Run("no managers", func(t *testing.T) {
		h := &localHandler{}
		h.Start()
	})
}

func TestLambdaHandlerInit(t *testing.T) {
	assert := assert.New(t)
	beforeType := os.Getenv(income.TypeEnv)

	t.Run("no errors", func(t *testing.T) {
		if err := os.Setenv(income.TypeEnv, income.GitHubType); err != nil {
			t.Fatal(err)
		}

		h := &lambdaHandler{}
		err := h.Init()
		assert.Nil(err)
		assert.IsType(h.in, &income.Manager{})
		assert.IsType(h.out, &outcome.Manager{})
	})

	t.Run("error", func(t *testing.T) {
		if err := os.Unsetenv(income.TypeEnv); err != nil {
			t.Fatal(err)
		}

		h := &lambdaHandler{}
		err := h.Init()
		assert.EqualError(err, fmt.Sprintf("Invalid %s", income.TypeEnv))
		assert.Nil(h.in)
		assert.Nil(h.out)
	})

	t.Cleanup(func(){
		if err := os.Setenv(income.TypeEnv, beforeType); err != nil {
			if err != nil {
				t.Fatal(err)
			}
		}
	})
}

func TestLambdaHandlerSrart(t *testing.T) {
	// TODO: lambda 依存のテストケースについて
}
