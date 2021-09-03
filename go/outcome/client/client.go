package client

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

const WebHookUrl = "SLACK_WEBHOOK_URL"

func NewClient() AbstractClient {
	// Only Slack
	return &SlackClient{}
}

type AbstractClient interface {
	Init() error
	Post(buff *bytes.Buffer) ([]byte, error)
}

type SlackClient struct {
	webHookUrl string
}

func (sc *SlackClient) Init() error {
	sc.webHookUrl = os.Getenv(WebHookUrl)
	if sc.webHookUrl == "" {
		return errors.New("WebhookUrl is brank")
	}
	return nil
}

func (sc *SlackClient) Post(msg *bytes.Buffer) ([]byte, error) {
	body, err := sc.request(msg)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (sc *SlackClient) request(msg *bytes.Buffer) ([]byte, error) {
	req, err := http.NewRequest("POST", sc.webHookUrl, msg)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	c := &http.Client{}
	resp, err := c.Do(req)
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
