package client

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/jarcoal/httpmock"
)

func TestNewClient(t *testing.T) {
	// Only Slack
	t.Run("Slack", func(t *testing.T) {
		c := NewClient()
		_, ok := c.(AbstractClient)
		assert.True(t, ok)
	})
}

func TestSlackClientInit(t *testing.T) {
	assert := assert.New(t)
	beforeWebHookUrl := os.Getenv(WebHookUrl)

	sc := &SlackClient{}

	t.Run("without WebHookUrl", func(t *testing.T) {
		if err := os.Unsetenv(WebHookUrl); err != nil {
			t.Fatal(err)
		}
		err := sc.Init()
		assert.EqualError(err, "WebhookUrl is brank")
	})

	t.Run("with WebHookUrl", func(t *testing.T) {
		mockUrl := "https://example.com"
		if err := os.Setenv(WebHookUrl, mockUrl); err != nil {
			t.Fatal(err)
		}
		err := sc.Init()
		assert.Nil(err)
		assert.Equal(sc.webHookUrl, mockUrl)
	})

	t.Cleanup(func(){
		if err := os.Setenv(WebHookUrl, beforeWebHookUrl); err != nil {
			t.Fatal(err)
		}
	})
}

func TestSlackClientPost(t *testing.T) {
	assert := assert.New(t)

	mockUrl := "https://example.com"
	msg := bytes.NewBufferString(`{"body": "test"}`)

	sc := &SlackClient{webHookUrl: mockUrl}

	t.Run("ok", func(t *testing.T) {
		resp := "ok"
		httpmock.RegisterResponder("POST", mockUrl,
			httpmock.NewStringResponder(200, resp))
		httpmock.Activate()

		body, err := sc.Post(msg)
		assert.Nil(err)
		assert.Equal(body, []byte(resp))

		t.Cleanup(func(){
			httpmock.DeactivateAndReset()
		})
	})

	t.Run("error", func(t *testing.T) {
		eerr := "mocked"
		httpmock.RegisterResponder("POST", mockUrl,
			func(req *http.Request) (*http.Response, error) {
				return httpmock.NewStringResponse(400, eerr), nil
			},
		)
		httpmock.Activate()

		_, err := sc.Post(msg)
		assert.EqualError(err, fmt.Sprintf("Status: 400, StatusCode 400, Body: %s", eerr))

		t.Cleanup(func(){
			httpmock.DeactivateAndReset()
		})
	})
}
