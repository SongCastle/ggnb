package message

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/google/go-github/v38/github"
)

const EventHeader = "x-github-event"

type fields = map[string]interface{}

type GitHubMessage struct {
	sender *string
	action string
	body *string
	target bool
}

func (gm *GitHubMessage) toGitHubEvent(payload interface{}) (interface{}, error) {
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

func (gm *GitHubMessage) Init(payload interface{}) error {
	event, err := gm.toGitHubEvent(payload)
	if err != nil {
		return err
	}
	// TODO: メッセージを充実する
	gm.target = true
	switch event := event.(type) {
	case *github.CommitCommentEvent:
		gm.sender = event.GetSender().Login
		gm.action = "コメントされました"
		gm.body = event.GetComment().Body
	case *github.CreateEvent:
		ref := event.GetRef()
		gm.sender = event.GetSender().Login
		gm.action = "ブランチ・タグが作成されました" // refType での判断も可能
		gm.body = &ref
	case *github.DeleteEvent:
		ref := event.GetRef()
		gm.sender = event.GetSender().Login
		gm.action = "ブランチ・タグが削除されました" // refType での判断も可能
		gm.body = &ref
	case *github.IssueCommentEvent:
		gm.sender = event.GetSender().Login
		gm.action = fmt.Sprintf("コメントが変更されました (%s)", event.GetAction())
		gm.body = event.GetIssue().URL
	case *github.IssuesEvent:
		gm.sender = event.GetSender().Login
		gm.action = fmt.Sprintf("PR / Issue (%s)", event.GetAction())
		gm.body = event.GetIssue().URL
	case *github.PullRequestEvent:
		gm.sender = event.GetSender().Login
		gm.action = fmt.Sprintf("PR (%s)", event.GetAction())
		gm.body = event.GetPullRequest().URL
	case *github.PullRequestReviewEvent:
		gm.sender = event.GetSender().Login
		gm.action = fmt.Sprintf("PR (%s)", event.GetAction())
		gm.body = event.GetPullRequest().URL
	case *github.PullRequestReviewCommentEvent: // PullRequestTargetEvent
		gm.sender = event.GetSender().Login
		gm.action = fmt.Sprintf("PR (%s)", event.GetAction())
		gm.body = event.GetPullRequest().URL
	case *github.PushEvent:
		ref := event.GetRef()
		gm.sender = event.GetSender().Login
		gm.action = "プッシュされました"
		gm.body = &ref
	default:
		gm.target = false
	}
	return nil
}

func (gm *GitHubMessage) NeedToDeliver() bool {
	return gm.target
}

func (gm *GitHubMessage) ToPayload() *bytes.Buffer {
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
								"short": true
							},
							{
								"title": "アクション",
								"value": "%s",
								"short": true
							},
							{
								"title": "内容",
								"value": "%s",
								"short": true
							}
						]
					}
				]
			}
			`, *gm.sender, gm.action, *gm.body,
		),
	)
}

func (gm *GitHubMessage) ToDummyPayload() *bytes.Buffer {
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
									"value": "bot",
									"short": true
								},
								{
									"title": "アクション",
									"value": "test",
									"short": true
								},
								{
									"title": "内容",
									"value": "good",
									"short": true
								}
							]
						}
					]
				}
			`,
		),
	)
}