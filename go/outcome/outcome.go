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
	Init(string)
	Post(buff *bytes.Buffer) error
}

type slackClient struct {
	webHookUrl string
}

func New() client {
	sc := slackClient{}
	sc.Init(os.Getenv(WebHookUrl))
	return &sc
}

func (sc *slackClient) Init(webHookUrl string) {
	sc.webHookUrl = webHookUrl
}

func (sc *slackClient) Post(buff *bytes.Buffer) error {
	if sc.webHookUrl == "" {
		return errors.New("WebhookUrl is brank")
	}

	req, err := http.NewRequest("POST", sc.webHookUrl, buff)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf(
		"Status: %s, StatusCode %d, Body: %s\n",
		resp.Status, resp.StatusCode, string(bodyBytes),
	)
	fmt.Println(msg)
	if resp.StatusCode != 200 {
		return errors.New(msg)
	}

	return nil
}
