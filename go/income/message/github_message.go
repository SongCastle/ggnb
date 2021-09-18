package message

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/SongCastle/ggnb/income/builder"
	"github.com/google/go-github/v38/github"
)

const (
	EventHeader = "x-github-event"
	EventHeaderCap = "X-GitHub-Event"
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

type fields = map[string]interface{}

type GitHubMessage struct {
	event interface{}
}

func (gm *GitHubMessage) Init(headers, body interface{}) error {
	_headers, ok := headers.(map[string]string)
	if !ok {
		return errors.New("invalid headers")
	}
	_body, ok := body.(*string)
	if !ok {
		return errors.New("invalid body")
	}
	event, err := gm.toGitHubEvent(_headers, _body)
	if err != nil {
		return err
	}
	gm.event = event
	return nil
}

func (gm *GitHubMessage) toGitHubEvent(headers map[string]string, body *string) (interface{}, error) {
	eventType, err := extractGitHubEvent(headers)
	if err != nil {
		return nil, err
	}
	return github.ParseWebHook(eventType, []byte(*body))
}

func extractGitHubEvent(headers map[string]string) (string, error) {
	eventType, ok := headers[EventHeader]
	if !ok {
		eventType, ok = headers[EventHeaderCap]
	}
	if !ok {
		return "", errors.New(fmt.Sprintf("missing %s header", EventHeader))
	}
	return eventType, nil
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

func (gm *GitHubMessage) ToDummyPayload() (*bytes.Buffer, error) {
	a := builder.NewAttachment()
	a.InsertField("アカウント", "bot", true)
	a.InsertField("アクション", "debug", true)
	a.InsertField("内容", "ok")
	return a.Build()
}

func buildCommitCommentEvent(e *commitCommentEvent) (*bytes.Buffer, error) {
	a := builder.NewAttachment()
	a.InsertField("アカウント", e.GetSender().GetLogin(), true)
	switch e.GetAction() {
	case "created":
		a.InsertField("アクション", "コメントされました", true)
		a.InsertField("コメント", e.GetComment().GetBody())
		a.InsertField("CommitID", e.GetComment().GetCommitID())
		a.InsertField("リンク", e.GetComment().GetHTMLURL())
	default:
		a.InsertField("アクション", fmt.Sprintf("CommitCommentEvent (%s)", e.GetAction()))
	}
	return a.Build()
}

func buildCreateEvent(e *createEvent) (*bytes.Buffer, error) {
	a := builder.NewAttachment()
	a.InsertField("アカウント", e.GetSender().GetLogin(), true)
	switch e.GetRefType() {
	case "branch":
		a.InsertField("アクション", "ブランチが作成されました", true)
		a.InsertField("ブランチ名", e.GetRef())
		a.InsertField("リンク", e.GetRepo().GetHTMLURL())
	case "tag":
		a.InsertField("アクション", "タグが作成されました", true)
		a.InsertField("タグ名", e.GetRef())
		a.InsertField("リンク", e.GetRepo().GetHTMLURL())
	}
	return a.Build()
}

func buildDeleteEvent(e *deleteEvent) (*bytes.Buffer, error) {
	a := builder.NewAttachment()
	a.InsertField("アカウント", e.GetSender().GetLogin(), true)
	switch e.GetRefType() {
	case "branch":
		a.InsertField("アクション", "ブランチが削除されました", true)
		a.InsertField("ブランチ名", e.GetRef())
		a.InsertField("リンク", e.GetRepo().GetHTMLURL())
	case "tag":
		a.InsertField("アクション", "タグが削除されました", true)
		a.InsertField("タグ名", e.GetRef())
		a.InsertField("リンク", e.GetRepo().GetHTMLURL())
	}
	return a.Build()
}

func buildIssueCommentEvent(e *issueCommentEvent) (*bytes.Buffer, error) {
	a := builder.NewAttachment()
	a.InsertField("アカウント", e.GetSender().GetLogin(), true)
	switch e.GetAction() {
	case "created":
		a.InsertField("アクション", "コメントされました", true)
		a.InsertField("コメント", e.GetComment().GetBody())
		a.InsertField("リンク", e.GetComment().GetHTMLURL())
	case "edited":
		a.InsertField("アクション", "コメントが変更されました", true)
		if body := e.GetChanges().GetBody(); body != nil {
			a.InsertField("コメント(変更前)", body.GetFrom())
			a.InsertField("コメント(変更後)", e.GetComment().GetBody())
		}
		a.InsertField("リンク", e.GetComment().GetHTMLURL())
	case "deleted":
		a.InsertField("アクション", "コメントが削除されました", true)
		a.InsertField("コメント", e.GetComment().GetBody())
		a.InsertField("リンク", e.GetComment().GetHTMLURL())
	default:
		a.InsertField("アクション", fmt.Sprintf("IssueCommentEvent (%s)", e.GetAction()))
	}
	return a.Build()
}

func buildIssuesEvent(e *issuesEvent) (*bytes.Buffer, error) {
	a := builder.NewAttachment()
	a.InsertField("アカウント", e.GetSender().GetLogin(), true)
	switch e.GetAction() {
	case "opened":
		a.InsertField("アクション", "Issue がオープンされました", true)
		a.InsertField("タイトル", e.GetIssue().GetTitle())
		a.InsertField("内容", e.GetIssue().GetBody())
		a.InsertField("リンク", e.GetIssue().GetHTMLURL())
	case "edited":
		a.InsertField("アクション", "Issue が編集されました", true)
		if title := e.GetChanges().GetTitle(); title != nil {
			a.InsertField("タイトル(変更前)", title.GetFrom())
			a.InsertField("タイトル(変更後)", e.GetIssue().GetTitle())
		}
		if body := e.GetChanges().GetBody(); body != nil {
			a.InsertField("内容(変更前)", body.GetFrom())
			a.InsertField("内容(変更後)", e.GetIssue().GetBody())
		}
		a.InsertField("リンク", e.GetIssue().GetHTMLURL())
	case "deleted":
		a.InsertField("アクション", "Issue が削除されました", true)
		a.InsertField("タイトル", e.GetIssue().GetTitle())
		a.InsertField("リンク", e.GetIssue().GetHTMLURL())
	case "closed":
		a.InsertField("アクション", "Issue がクローズされました", true)
		a.InsertField("タイトル", e.GetIssue().GetTitle())
		a.InsertField("リンク", e.GetIssue().GetHTMLURL())
	case "reopened":
		a.InsertField("アクション", "Issue が再オープンされました", true)
		a.InsertField("タイトル", e.GetIssue().GetTitle())
		a.InsertField("内容", e.GetIssue().GetBody())
		a.InsertField("リンク", e.GetIssue().GetHTMLURL())
	case "assigned":
		a.InsertField("アクション", "Issue にアサインされました", true)
		a.InsertField("対象者", e.GetAssignee().GetLogin())
		a.InsertField("リンク", e.GetIssue().GetHTMLURL())
	case "unassigned":
		a.InsertField("アクション", "Issue にアンアサインされました", true)
		a.InsertField("対象者", e.GetAssignee().GetLogin())
		a.InsertField("リンク", e.GetIssue().GetHTMLURL())
	case "labeled":
		a.InsertField("アクション", "Issue にラベルが付与されました", true)
		a.InsertField("ラベル", e.GetLabel().GetName())
		a.InsertField("リンク", e.GetIssue().GetHTMLURL())
	case "unlabeled":
		a.InsertField("アクション", "Issue のラベルが外されました", true)
		a.InsertField("ラベル", e.GetLabel().GetName())
		a.InsertField("リンク", e.GetIssue().GetHTMLURL())
	case "locked":
		a.InsertField("アクション", "Issue がロックされました", true)
		a.InsertField("リンク", e.GetIssue().GetHTMLURL())
	case "unlocked":
		a.InsertField("アクション", "Issue のロックが解除されました", true)
		a.InsertField("リンク", e.GetIssue().GetHTMLURL())
	case "pinned":
		a.InsertField("アクション", "Issue がピン留めされました", true)
		a.InsertField("リンク", e.GetIssue().GetHTMLURL())
	case "unpinned":
		a.InsertField("アクション", "Issue のピン留めが解除されました", true)
		a.InsertField("リンク", e.GetIssue().GetHTMLURL())
	case "transferred": // TODO: 譲渡先の URL を取得する
		a.InsertField("アクション", "Issue が譲渡されました", true)
		a.InsertField("リンク(譲渡前)", e.GetIssue().GetHTMLURL())
	case "milestoned":
		a.InsertField("アクション", "マイルストーンが設定されました", true)
		a.InsertField("マイルストーン", e.GetIssue().GetMilestone().GetTitle())
		a.InsertField("リンク", e.GetIssue().GetHTMLURL())
	case "demilestoned":
		a.InsertField("アクション", "マイルストーンが解除されました", true)
		a.InsertField("リンク", e.GetIssue().GetHTMLURL())
	default:
		a.InsertField("アクション", fmt.Sprintf("IssuesEvent (%s)", e.GetAction()))
	}
	return a.Build()
}

func buildPullRequestEvent(e *pullRequestEvent) (*bytes.Buffer, error) {
	a := builder.NewAttachment()
	a.InsertField("アカウント", e.GetSender().GetLogin(), true)
	switch e.GetAction() {
	case "opened":
		a.InsertField("アクション", "PR がオープンされました", true)
		a.InsertField("タイトル", e.GetPullRequest().GetTitle())
		a.InsertField("内容", e.GetPullRequest().GetBody())
		a.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "edited":
		a.InsertField("アクション", "PR が編集されました", true)
		if title := e.GetChanges().GetTitle(); title != nil {
			a.InsertField("タイトル(変更前)", title.GetFrom())
			a.InsertField("タイトル(変更後)", e.GetPullRequest().GetTitle())
		}
		if body := e.GetChanges().GetBody(); body != nil {
			a.InsertField("内容(変更前)", body.GetFrom())
			a.InsertField("内容(変更後)", e.GetPullRequest().GetBody())
		}
		a.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "closed":
		a.InsertField("アクション", "PR がクローズされました", true)
		a.InsertField("タイトル", e.GetPullRequest().GetTitle())
		a.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "reopened":
		a.InsertField("アクション", "PR が再オープンされました", true)
		a.InsertField("タイトル", e.GetPullRequest().GetTitle())
		a.InsertField("内容", e.GetPullRequest().GetBody())
		a.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "assigned":
		a.InsertField("アクション", "PR にアサインされました", true)
		a.InsertField("対象者", e.GetAssignee().GetLogin())
		a.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "unassigned":
		a.InsertField("アクション", "PR にアンアサインされました", true)
		a.InsertField("対象者", e.GetAssignee().GetLogin())
		a.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "review_requested":
		a.InsertField("アクション", "PR のレビューをお願いされました", true)
		a.InsertField("対象者", e.GetRequestedReviewer().GetLogin())
		a.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "review_request_removed":
		a.InsertField("アクション", "PR のレビュー要求が取下げされました", true)
		a.InsertField("対象者", e.GetRequestedReviewer().GetLogin())
		a.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "ready_for_review":
		a.InsertField("アクション", "PR の準備が整いました", true)
		a.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "labeled":
		a.InsertField("アクション", "PR にラベルが付与されました", true)
		a.InsertField("ラベル", e.GetLabel().GetName())
		a.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "unlabeled":
		a.InsertField("アクション", "PR のラベルが外されました", true)
		a.InsertField("ラベル", e.GetLabel().GetName())
		a.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "locked":
		a.InsertField("アクション", "PR がロックされました", true)
		a.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "unlocked":
		a.InsertField("アクション", "PR のロックが解除されました", true)
		a.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "synchronize":
		a.InsertField("アクション", "PullRequestEvent synchronize")
	default:
		a.InsertField("アクション", fmt.Sprintf("PullRequestEvent (%s)", e.GetAction()))
	}
	return a.Build()
}

func buildPullRequestReviewEvent(e *pullRequestReviewEvent) (*bytes.Buffer, error) {
	a := builder.NewAttachment()
	a.InsertField("アカウント", e.GetSender().GetLogin(), true)
	switch e.GetAction() {
	case "submitted":
		a.InsertField("アクション", "PR のレビューがされました", true)
		a.InsertField("タイトル", e.GetPullRequest().GetTitle())
		a.InsertField("内容", e.GetReview().GetBody())
		a.InsertField("リンク", e.GetReview().GetHTMLURL())
	default:
		a.InsertField("アクション", fmt.Sprintf("PullRequestReviewEvent (%s)", e.GetAction()))
	}
	return a.Build()
}

func buildPullRequestReviewCommentEvent(e *pullRequestReviewCommentEvent) (*bytes.Buffer, error) {
	a := builder.NewAttachment()
	a.InsertField("アカウント", e.GetSender().GetLogin(), true)
	switch e.GetAction() {
	case "created":
		a.InsertField("アクション", "PR にコメントされました", true)
		a.InsertField("コメント", e.GetComment().GetBody())
		a.InsertField("リンク", e.GetComment().GetHTMLURL())
	case "edited":
		a.InsertField("アクション", "PR のコメントが変更されました", true)
		if body := e.GetChanges().GetBody(); body != nil {
			a.InsertField("コメント(変更前)", e.GetChanges().GetBody().GetFrom())
			a.InsertField("コメント(変更後)", e.GetComment().GetBody())
		}
		a.InsertField("リンク", e.GetComment().GetHTMLURL())
	case "deleted":
		a.InsertField("アクション", "PR のコメントが削除されました", true)
		a.InsertField("コメント", e.GetComment().GetBody())
		a.InsertField("リンク", e.GetComment().GetHTMLURL())
	default:
		a.InsertField("アクション", fmt.Sprintf("PullRequestReviewCommentEvent (%s)", e.GetAction()))
	}
	return a.Build()
}

func buildPullRequestTargetEvent(e *pullRequestTargetEvent) (*bytes.Buffer, error) {
	a := builder.NewAttachment()
	a.InsertField("アカウント", e.GetSender().GetLogin(), true)
	switch e.GetAction() {
	case "opened":
		a.InsertField("アクション", "PR がオープンされました", true)
		a.InsertField("タイトル", e.GetPullRequest().GetTitle())
		a.InsertField("内容", e.GetPullRequest().GetBody())
		a.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "edited":
		a.InsertField("アクション", "PR が編集されました", true)
		if title := e.GetChanges().GetTitle(); title != nil {
			a.InsertField("タイトル(変更前)", title.GetFrom())
			a.InsertField("タイトル(変更後)", e.GetPullRequest().GetTitle())
		}
		if body := e.GetChanges().GetBody(); body != nil {
			a.InsertField("内容(変更前)", body.GetFrom())
			a.InsertField("内容(変更後)", e.GetPullRequest().GetBody())
		}
		a.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "closed":
		a.InsertField("アクション", "PR がクローズされました", true)
		a.InsertField("タイトル", e.GetPullRequest().GetTitle())
		a.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "reopened":
		a.InsertField("アクション", "PR が再オープンされました", true)
		a.InsertField("タイトル", e.GetPullRequest().GetTitle())
		a.InsertField("内容", e.GetPullRequest().GetBody())
		a.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "assigned":
		a.InsertField("アクション", "PR にアサインされました", true)
		a.InsertField("対象者", e.GetAssignee().GetLogin())
		a.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "unassigned":
		a.InsertField("アクション", "PR にアンアサインされました", true)
		a.InsertField("対象者", e.GetAssignee().GetLogin())
		a.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "review_requested":
		a.InsertField("アクション", "PR のレビューをお願いされました", true)
		a.InsertField("対象者", e.GetRequestedReviewer().GetLogin())
		a.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "review_request_removed":
		a.InsertField("アクション", "PR のレビュー要求が取下げされました", true)
		a.InsertField("対象者", e.GetRequestedReviewer().GetLogin())
		a.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "ready_for_review":
		a.InsertField("アクション", "PR の準備が整いました", true)
		a.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "labeled":
		a.InsertField("アクション", "PR にラベルが付与されました", true)
		a.InsertField("ラベル", e.GetLabel().GetName())
		a.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "unlabeled":
		a.InsertField("アクション", "PR のラベルが外されました", true)
		a.InsertField("ラベル", e.GetLabel().GetName())
		a.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "locked":
		a.InsertField("アクション", "PR がロックされました", true)
		a.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "unlocked":
		a.InsertField("アクション", "PR のロックが解除されました", true)
		a.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "synchronize":
		a.InsertField("アクション", "PullRequestEvent synchronize")
	default:
		a.InsertField("アクション", fmt.Sprintf("PullRequestTargetEvent (%s)", e.GetAction()))
	}
	return a.Build()
}

func buildPushEvent(e *pushEvent) (*bytes.Buffer, error) {
	a := builder.NewAttachment()
	a.InsertField("アカウント", e.GetSender().GetLogin(), true)
	a.InsertField("アクション", "プッシュされました", true)
	a.InsertField("対象", e.GetRef())
	if e.Commits != nil {
		var b strings.Builder
		for _, c := range e.Commits {
			b.WriteString(
				fmt.Sprintf("<%s|%s> %s\n", c.GetURL(), c.GetID()[:7], c.GetMessage()),
			)
		}
		a.InsertField("Commit", b.String())
	}
	a.InsertField("リンク", e.GetRepo().GetHTMLURL())
	return a.Build()
}
