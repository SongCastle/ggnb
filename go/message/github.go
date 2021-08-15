package message

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/google/go-github/v38/github"
)

type commitCommentEvent = github.CommitCommentEvent
type createEvent = github.CreateEvent
type deleteEvent = github.DeleteEvent
type issueCommentEvent = github.IssueCommentEvent
type issuesEvent = github.IssuesEvent
type pullRequestEvent = github.PullRequestEvent
type pullRequestReviewEvent = github.PullRequestReviewEvent
type pullRequestReviewCommentEvent = github.PullRequestReviewCommentEvent
type pullRequestTargetEvent = github.PullRequestTargetEvent
type pushEvent = github.PushEvent

const EventHeader = "x-github-event"

type fields = map[string]interface{}

type GitHubMessage struct {
	event interface{}
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
	gm.event = event
	return nil
}

func (gm *GitHubMessage) ToPayload() (*bytes.Buffer, error) {
	switch event := gm.event.(type) {
	case *commitCommentEvent:
		return buildCommitCommentEvent(event)
	case *createEvent:
		return buildCreateEvent(event)
	case *deleteEvent:
		return buildDeleteEvent(event)
	case *issueCommentEvent:
		return buildIssueCommentEvent(event)
	case *issuesEvent:
		return buildIssuesEvent(event)
	case *pullRequestEvent:
		return buildPullRequestEvent(event)
	case *pullRequestReviewEvent:
		return buildPullRequestReviewEvent(event)
	case *pullRequestReviewCommentEvent:
		return buildPullRequestReviewCommentEvent(event)
	case *pullRequestTargetEvent:
		return buildPullRequestTargetEvent(event)
	case *pushEvent:
		return buildPushEvent(event)
	default:
		return nil, nil
	}
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
