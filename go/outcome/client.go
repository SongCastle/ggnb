package outcome

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

const WebHookUrl = "SLACK_WEBHOOK_URL"

type client interface {
	Init()
	Post(buff *bytes.Buffer) ([]byte, error)
}

type slackClient struct {
	webHookUrl string
}

func (sc *slackClient) Init() {
	sc.webHookUrl = os.Getenv(WebHookUrl)
}

func (sc *slackClient) Post(msg *bytes.Buffer) ([]byte, error) {
	if sc.webHookUrl == "" {
		return nil, errors.New("WebhookUrl is brank")
	}
	body, err := sc.request(msg)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (sc *slackClient) request(msg *bytes.Buffer) ([]byte, error) {
	req, err := http.NewRequest("POST", sc.webHookUrl, msg)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := fmt.Sprintf(
		"Status: %s, StatusCode %d, Body: %s",
		resp.Status, resp.StatusCode, string(body),
	)
	fmt.Println(result)

	if resp.StatusCode != 200 {
		return nil, errors.New(result)
	}
	return body, nil
}
