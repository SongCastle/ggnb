package handler

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/SongCastle/ggnb/income"
	"github.com/SongCastle/ggnb/income/message"
	"github.com/SongCastle/ggnb/outcome"
	"github.com/SongCastle/ggnb/outcome/client"
	"github.com/stretchr/testify/assert"
)

func TestHandlerNew(t *testing.T) {
	assert := assert.New(t)
	beforeDebug := os.Getenv(DEBUG)

	m := &message.MockedMessage{}
	c := &client.MockedClient{}

	t.Run("local", func(t *testing.T) {
		if err := os.Setenv(DEBUG, "1"); err != nil {
			t.Fatal(err)
		}
		h := New(m, c)
		assert.IsType(h, &localHandler{})
	})

	t.Run("lambda", func(t *testing.T) {
		if err := os.Unsetenv(DEBUG); err != nil {
			t.Fatal(err)
		}
		h := New(m, c)
		assert.IsType(h, &lambdaHandler{})
	})

	t.Cleanup(func(){
		if err := os.Setenv(DEBUG, beforeDebug); err != nil {
			t.Fatal(err)
		}
	})
}

func TestLocalHandlerInit(t *testing.T) {
	h := &localHandler{}
	assert.NotPanics(
		t,
		func() {
			h.Init(
				&income.MockedIncomeManager{},
				&outcome.MockedOutcomeManager{},
			)
		},
	)
}

func TestLocalHandlerStart(t *testing.T) {
	assert := assert.New(t)

	t.Run("no errors", func(t *testing.T) {
		msg := bytes.NewBufferString(`{"body": "test"}`)
		in := &income.MockedIncomeManager{}
		in.On("BuildDummyMessage").Return(msg, nil)

		out := &outcome.MockedOutcomeManager{}
		out.On("Send", msg).Return(nil)
		out.On("ReportErrorIf", nil).Return(nil)

		h := &localHandler{In: in, Out: out}
		assert.NotPanics(func() { h.Start() })
	})

	t.Run("error", func(t *testing.T) {
		var b *bytes.Buffer
		err := errors.New("mocked")

		in := &income.MockedIncomeManager{}
		in.On("BuildDummyMessage").Return(b, err)

		out := &outcome.MockedOutcomeManager{}
		out.On("ReportErrorIf", err).Return(err)

		h := &localHandler{In: in, Out: out}
		assert.NotPanics(func() { h.Start() })
	})
}

func TestLambdaHandlerInit(t *testing.T) {
	h := &lambdaHandler{}
	assert.NotPanics(
		t,
		func() {
			h.Init(
				&income.MockedIncomeManager{},
				&outcome.MockedOutcomeManager{},
			)
		},
	)
}

func TestLambdaHandlerSrart(t *testing.T) {
	// TODO: lambda 依存のテストケースについて
}
