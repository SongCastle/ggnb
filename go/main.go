package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/go-github/v38/github"
)

const (
	WebHookUrl = "SLACK_WEBHOOK_URL"
	DEBUG = "DEBUG"
	EventHeader = "x-github-event"
)

type fields = map[string]interface{}

type slackClient struct {
	webHookUrl string
}

type slackResponse struct {
	Status string
	StatusCode int
	Body string
}

type slackMessage struct {
	Sender string `json:"sender,omitempty"`
	Action string `json:"action,omitempty"`
	Body string `json:"body,omitempty"`
	Built bool
}

func (sm *slackMessage) toWebHookPayload(payload interface{}) (interface{}, error) {
	_payload := payload.(fields)
	// check header
	headers, ok := _payload["headers"]
	if !ok {
		return nil, errors.New(fmt.Sprintf("missing headers"))
	}
	eventType, ok := headers.(fields)[EventHeader]
	if !ok {
		return nil, errors.New(fmt.Sprintf("missing %s header", EventHeader))
	}
	// check body
	body, ok := _payload["body"]
	if !ok {
		return nil, errors.New(fmt.Sprintf("missing body"))
	}
	return github.ParseWebHook(eventType.(string), []byte(body.(string)))
}

func (sm *slackMessage) Build(payload interface{}) error {
	event, err := sm.toWebHookPayload(payload)
	if err != nil {
		return err
	}
	// TODO: メッセージを充実する
	sm.Built = true
	switch event := event.(type) {
	case *github.CommitCommentEvent:
		sm.Sender = *event.GetSender().Login
		sm.Action = "コメントされました"
		sm.Body = *event.GetComment().Body
	case *github.CreateEvent:
		sm.Sender =  *event.GetSender().Login
		sm.Action = "ブランチ・タグが作成されました" // refType での判断も可能
		sm.Body = event.GetRef()
	case *github.DeleteEvent:
		sm.Sender = *event.GetSender().Login
		sm.Action = "ブランチ・タグが削除されました" // refType での判断も可能
		sm.Body = event.GetRef()
	case *github.IssueCommentEvent:
		sm.Sender = *event.GetSender().Login
		sm.Action = fmt.Sprintf("コメントが変更されました (%s)", event.GetAction())
		sm.Body = *event.GetIssue().URL
	case *github.IssuesEvent:
		sm.Sender = *event.GetSender().Login
		sm.Action = fmt.Sprintf("PR / Issue (%s)", event.GetAction())
		sm.Body = *event.GetIssue().URL
	case *github.PullRequestEvent:
		sm.Sender = *event.GetSender().Login
		sm.Action = fmt.Sprintf("PR (%s)", event.GetAction())
		sm.Body = *event.GetPullRequest().URL
	case *github.PullRequestReviewEvent:
		sm.Sender = *event.GetSender().Login
		sm.Action = fmt.Sprintf("PR (%s)", event.GetAction())
		sm.Body = *event.GetPullRequest().URL
	case *github.PullRequestReviewCommentEvent: // PullRequestTargetEvent
		sm.Sender = *event.GetSender().Login
		sm.Action = fmt.Sprintf("PR (%s)", event.GetAction())
		sm.Body = *event.GetPullRequest().URL
	case *github.PushEvent:
		sm.Sender = *event.GetSender().Login
		sm.Action = "プッシュされました"
		sm.Body = event.GetRef()
	default:
		sm.Built = false
	}
	return nil
}

func (sm *slackMessage) ToPayload() *bytes.Buffer {
	if !sm.Built {
		return nil
	}
	return bytes.NewBufferString(
		fmt.Sprintf(
			`
			{
				"attachments": [
					{
						"fallback": "GitHub Notifitation",
						"color": "#2eb886",
						"title": "GitHub Notification",
						"title_link": "https://api.slack.com/",
						"fields": [
							{
								"title": "アカウント",
								"value": "%s",
								"short": "true"
							},
							{
								"title": "アクション",
								"value": "%s",
								"short": "true"
							},
							{
								"title": "内容",
								"value": "%s",
								"short": "true"
							}
						]
					}
				]
			}
			`, sm.Sender, sm.Action, sm.Body,
		),
	)
}

func (sc *slackClient) Init(webHookUrl string) {
	sc.webHookUrl = webHookUrl
}

func (sc *slackClient) Post(buff *bytes.Buffer) (*slackResponse, error) {
	if sc.webHookUrl == "" {
		return nil, errors.New("slack webhook url is brank")
	}

	req, err := http.NewRequest("POST", sc.webHookUrl, buff)
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

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	sresp := slackResponse{
		Status: resp.Status,
		StatusCode: resp.StatusCode,
		Body: string(bodyBytes),
	}
	return &sresp, nil
}

func postToSlack(payload interface{}) error {
	sm := &slackMessage{}
	if err := sm.Build(payload); err != nil {
		fmt.Printf("failed: %v\n", err)
		return err
	}
	msg := sm.ToPayload()
	if msg == nil {
		fmt.Println("skipped")
		return nil
	}

	sc := slackClient{}
	sc.Init(os.Getenv(WebHookUrl))
	sresp, err := sc.Post(msg)
	if err != nil {
		fmt.Printf("failed: %v\n", err)
		return err
	}
	if (*sresp).StatusCode == 200 {
		fmt.Printf("succeeded: %v\n", *sresp)
		return nil
	}
	err_msg := fmt.Sprintf("failed: %v", *sresp)
	fmt.Println(err_msg)
	return errors.New(err_msg)
}

func HandleRequest(ctx context.Context, payload interface{}) error {
	return postToSlack(payload)
}

func onLambda() bool {
	return os.Getenv(DEBUG) == ""
}

func main() {
	if onLambda() {
		lambda.Start(HandleRequest)
		return
	}
	body := `{"sender": {"login": "bot"}, "comment": {"body": "good"}}`
	headers := map[string]interface{}{EventHeader: "commit_comment"}
	postToSlack(
		map[string]interface{}{"body": body, "headers": headers},
	)
}
