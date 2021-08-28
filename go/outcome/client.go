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
	Post(buff *bytes.Buffer) error
}

type slackClient struct {
	webHookUrl string
}

func (sc *slackClient) Init() {
	sc.webHookUrl = os.Getenv(WebHookUrl)
}

func (sc *slackClient) Post(msg *bytes.Buffer) error {
	if sc.webHookUrl == "" {
		return errors.New("WebhookUrl is brank")
	}

	req, err := http.NewRequest("POST", sc.webHookUrl, msg)
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

	result := fmt.Sprintf(
		"Status: %s, StatusCode %d, Body: %s",
		resp.Status, resp.StatusCode, string(bodyBytes),
	)

	fmt.Println(result)
	if resp.StatusCode != 200 {
		return errors.New(result)
	}
	return nil
}
