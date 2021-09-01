package outcome

import (
	"bytes"
	"net/http"
	"os"
	"testing"

	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/jarcoal/httpmock"
)

func TestSlackClientInitAndPost(t *testing.T) {
	assert := assert.New(t)

	mockUrl := "https://example.com"
	msg := bytes.NewBufferString(`{"body": "test"}`)

	t.Run("without WebHookUrl", func(t *testing.T) {
		sc := &slackClient{}
		_, err := sc.Post(msg)
		assert.EqualError(err, "WebhookUrl is brank")
	})

	t.Run("with WebHookUrl", func(t *testing.T) {
		beforeWebHookUrl := os.Getenv(WebHookUrl)
		if err := os.Setenv(WebHookUrl, mockUrl); err != nil {
			if err != nil {
				t.Fatal(err)
			}
		}

		sc := &slackClient{}
		sc.Init()

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
			merr := "mocked"
			httpmock.RegisterResponder("POST", mockUrl,
				func(req *http.Request) (*http.Response, error) {
					return httpmock.NewStringResponse(400, merr), nil
				},
			)
			httpmock.Activate()

			_, err := sc.Post(msg)
			assert.EqualError(err, fmt.Sprintf("Status: 400, StatusCode 400, Body: %s", merr))

			t.Cleanup(func(){
				httpmock.DeactivateAndReset()
			})
		})

		t.Cleanup(func(){
			if err := os.Setenv(WebHookUrl, beforeWebHookUrl); err != nil {
				if err != nil {
					t.Fatal(err)
				}
			}
		})
	})
}
