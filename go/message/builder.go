package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type payloadField struct {
	Title *string `json:"title"`
	Value *string `json:"value"`
	Short bool `json:"short"`
}

func (pf *payloadField) Build() ([]byte, error) {
	return json.Marshal(pf)
}

type payloadBuilder struct {
	fields []*payloadField
}

func getShort(short ...bool) bool {
	if len(short) == 0 {
		return false	
	}
	return short[0]
}

func (pb *payloadBuilder) InsertField(title, value string, short ...bool) {
	pb.fields = append(
		pb.fields,
		&payloadField{Title: &title, Value: &value, Short: getShort(short...)},
	)
}

func (pb *payloadBuilder) Build() (*bytes.Buffer, error) {
	var sb strings.Builder
	for _, f := range pb.fields {
		b, err := f.Build()
		if err != nil {
			return nil, err
		}
		if _, err := sb.Write(b); err != nil {
			return nil, err
		}
		if _, err := sb.WriteString(","); err != nil {
			return nil, err
		}
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
						"title_link": "https://github.com/SongCastle/ggnb",
						"fields": [%s]
					}
				]
			}
			`, sb.String(),
		),
	), nil
}

func buildCommitCommentEvent(e *commitCommentEvent) (*bytes.Buffer, error) {
	pb := payloadBuilder{}
	pb.InsertField("アカウント", e.GetSender().GetLogin(), true)
	switch e.GetAction() {
	case "created":
		pb.InsertField("アクション", "コメントされました", true)
		pb.InsertField("コメント", e.GetComment().GetBody())
		pb.InsertField("CommitID", e.GetComment().GetCommitID())
		pb.InsertField("リンク", e.GetComment().GetHTMLURL())
	default:
		pb.InsertField("アクション", fmt.Sprintf("CommitCommentEvent (%s)", e.GetAction()))
	}
	return pb.Build()
}

func buildCreateEvent(e *createEvent) (*bytes.Buffer, error) {
	pb := payloadBuilder{}
	pb.InsertField("アカウント", e.GetSender().GetLogin(), true)
	switch e.GetRefType() {
	case "branch":
		pb.InsertField("アクション", "ブランチが作成されました", true)
		pb.InsertField("ブランチ名", e.GetRef())
		pb.InsertField("リンク", e.GetRepo().GetHTMLURL())
	case "tag":
		pb.InsertField("アクション", "タグが作成されました", true)
		pb.InsertField("タグ名", e.GetRef())
		pb.InsertField("リンク", e.GetRepo().GetHTMLURL())
	}
	return pb.Build()
}

func buildDeleteEvent(e *deleteEvent) (*bytes.Buffer, error) {
	pb := payloadBuilder{}
	pb.InsertField("アカウント", e.GetSender().GetLogin(), true)
	switch e.GetRefType() {
	case "branch":
		pb.InsertField("アクション", "ブランチが削除されました", true)
		pb.InsertField("ブランチ名", e.GetRef())
		pb.InsertField("リンク", e.GetRepo().GetHTMLURL())
	case "tag":
		pb.InsertField("アクション", "タグが削除されました", true)
		pb.InsertField("タグ名", e.GetRef())
		pb.InsertField("リンク", e.GetRepo().GetHTMLURL())
	}
	return pb.Build()
}

func buildIssueCommentEvent(e *issueCommentEvent) (*bytes.Buffer, error) {
	pb := payloadBuilder{}
	pb.InsertField("アカウント", e.GetSender().GetLogin(), true)
	switch e.GetAction() {
	case "created":
		pb.InsertField("アクション", "コメントされました", true)
		pb.InsertField("コメント", e.GetComment().GetBody())
		pb.InsertField("リンク", e.GetComment().GetHTMLURL())
	case "edited":
		pb.InsertField("アクション", "コメントが変更されました", true)
		if body := e.GetChanges().GetBody(); body != nil {
			pb.InsertField("コメント(変更前)", body.GetFrom())
			pb.InsertField("コメント(変更後)", e.GetComment().GetBody())
		}
		pb.InsertField("リンク", e.GetComment().GetHTMLURL())
	case "deleted":
		pb.InsertField("アクション", "コメントが削除されました", true)
		pb.InsertField("コメント", e.GetComment().GetBody())
		pb.InsertField("リンク", e.GetComment().GetHTMLURL())
	default:
		pb.InsertField("アクション", fmt.Sprintf("IssueCommentEvent (%s)", e.GetAction()))
	}
	return pb.Build()
}

func buildIssuesEvent(e *issuesEvent) (*bytes.Buffer, error) {
	pb := payloadBuilder{}
	pb.InsertField("アカウント", e.GetSender().GetLogin(), true)
	switch e.GetAction() {
	case "opened":
		pb.InsertField("アクション", "Issue がオープンされました", true)
		pb.InsertField("タイトル", e.GetIssue().GetTitle())
		pb.InsertField("内容", e.GetIssue().GetBody())
		pb.InsertField("リンク", e.GetIssue().GetHTMLURL())
	case "edited":
		pb.InsertField("アクション", "Issue が編集されました", true)
		if title := e.GetChanges().GetTitle(); title != nil {
			pb.InsertField("タイトル(変更前)", title.GetFrom())
			pb.InsertField("タイトル(変更後)", e.GetIssue().GetTitle())
		}
		if body := e.GetChanges().GetBody(); body != nil {
			pb.InsertField("内容(変更前)", body.GetFrom())
			pb.InsertField("内容(変更後)", e.GetIssue().GetBody())
		}
		pb.InsertField("リンク", e.GetIssue().GetHTMLURL())
	case "deleted":
		pb.InsertField("アクション", "Issue が削除されました", true)
		pb.InsertField("タイトル", e.GetIssue().GetTitle())
		pb.InsertField("リンク", e.GetIssue().GetHTMLURL())
	case "closed":
		pb.InsertField("アクション", "Issue がクローズされました", true)
		pb.InsertField("タイトル", e.GetIssue().GetTitle())
		pb.InsertField("リンク", e.GetIssue().GetHTMLURL())
	case "reopened":
		pb.InsertField("アクション", "Issue が再オープンされました", true)
		pb.InsertField("タイトル", e.GetIssue().GetTitle())
		pb.InsertField("内容", e.GetIssue().GetBody())
		pb.InsertField("リンク", e.GetIssue().GetHTMLURL())
	case "assigned":
		pb.InsertField("アクション", "Issue にアサインされました", true)
		pb.InsertField("対象者", e.GetAssignee().GetLogin())
		pb.InsertField("リンク", e.GetIssue().GetHTMLURL())
	case "unassigned":
		pb.InsertField("アクション", "Issue にアンアサインされました", true)
		pb.InsertField("対象者", e.GetAssignee().GetLogin())
		pb.InsertField("リンク", e.GetIssue().GetHTMLURL())
	case "labeled":
		pb.InsertField("アクション", "Issue にラベルが付与されました", true)
		pb.InsertField("ラベル", e.GetLabel().GetName())
		pb.InsertField("リンク", e.GetIssue().GetHTMLURL())
	case "unlabeled":
		pb.InsertField("アクション", "Issue のラベルが外されました", true)
		pb.InsertField("ラベル", e.GetLabel().GetName())
		pb.InsertField("リンク", e.GetIssue().GetHTMLURL())
	case "locked":
		pb.InsertField("アクション", "Issue がロックされました", true)
		pb.InsertField("リンク", e.GetIssue().GetHTMLURL())
	case "unlocked":
		pb.InsertField("アクション", "Issue のロックが解除されました", true)
		pb.InsertField("リンク", e.GetIssue().GetHTMLURL())
	case "pinned":
		pb.InsertField("アクション", "Issue がピン留めされました", true)
		pb.InsertField("リンク", e.GetIssue().GetHTMLURL())
	case "unpinned":
		pb.InsertField("アクション", "Issue のピン留めが解除されました", true)
		pb.InsertField("リンク", e.GetIssue().GetHTMLURL())
	case "transferred": // TODO: 譲渡先の URL を取得する
		pb.InsertField("アクション", "Issue が譲渡されました", true)
		pb.InsertField("リンク(譲渡前)", e.GetIssue().GetHTMLURL())
	case "milestoned":
		pb.InsertField("アクション", "マイルストーンが設定されました", true)
		pb.InsertField("マイルストーン", e.GetIssue().GetMilestone().GetTitle())
		pb.InsertField("リンク", e.GetIssue().GetHTMLURL())
	case "demilestoned":
		pb.InsertField("アクション", "マイルストーンが解除されました", true)
		pb.InsertField("リンク", e.GetIssue().GetHTMLURL())
	default:
		pb.InsertField("アクション", fmt.Sprintf("IssuesEvent (%s)", e.GetAction()))
	}
	return pb.Build()
}

func buildPullRequestEvent(e *pullRequestEvent) (*bytes.Buffer, error) {
	pb := payloadBuilder{}
	pb.InsertField("アカウント", e.GetSender().GetLogin(), true)
	switch e.GetAction() {
	case "opened":
		pb.InsertField("アクション", "PR がオープンされました", true)
		pb.InsertField("タイトル", e.GetPullRequest().GetTitle())
		pb.InsertField("内容", e.GetPullRequest().GetBody())
		pb.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "edited":
		pb.InsertField("アクション", "PR が編集されました", true)
		if title := e.GetChanges().GetTitle(); title != nil {
			pb.InsertField("タイトル(変更前)", title.GetFrom())
			pb.InsertField("タイトル(変更後)", e.GetPullRequest().GetTitle())
		}
		if body := e.GetChanges().GetBody(); body != nil {
			pb.InsertField("内容(変更前)", body.GetFrom())
			pb.InsertField("内容(変更後)", e.GetPullRequest().GetBody())
		}
		pb.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "closed":
		pb.InsertField("アクション", "PR がクローズされました", true)
		pb.InsertField("タイトル", e.GetPullRequest().GetTitle())
		pb.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "reopened":
		pb.InsertField("アクション", "PR が再オープンされました", true)
		pb.InsertField("タイトル", e.GetPullRequest().GetTitle())
		pb.InsertField("内容", e.GetPullRequest().GetBody())
		pb.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "assigned":
		pb.InsertField("アクション", "PR にアサインされました", true)
		pb.InsertField("対象者", e.GetAssignee().GetLogin())
		pb.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "unassigned":
		pb.InsertField("アクション", "PR にアンアサインされました", true)
		pb.InsertField("対象者", e.GetAssignee().GetLogin())
		pb.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "review_requested":
		pb.InsertField("アクション", "PR のレビューをお願いされました", true)
		pb.InsertField("対象者", e.GetRequestedReviewer().GetLogin())
		pb.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "review_request_removed":
		pb.InsertField("アクション", "PR のレビュー要求が取下げされました", true)
		pb.InsertField("対象者", e.GetRequestedReviewer().GetLogin())
		pb.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "ready_for_review":
		pb.InsertField("アクション", "PR の準備が整いました", true)
		pb.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "labeled":
		pb.InsertField("アクション", "PR にラベルが付与されました", true)
		pb.InsertField("ラベル", e.GetLabel().GetName())
		pb.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "unlabeled":
		pb.InsertField("アクション", "PR のラベルが外されました", true)
		pb.InsertField("ラベル", e.GetLabel().GetName())
		pb.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "locked":
		pb.InsertField("アクション", "PR がロックされました", true)
		pb.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "unlocked":
		pb.InsertField("アクション", "PR のロックが解除されました", true)
		pb.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "synchronize":
		pb.InsertField("アクション", "PullRequestEvent synchronize")
	default:
		pb.InsertField("アクション", fmt.Sprintf("PullRequestEvent (%s)", e.GetAction()))
	}
	return pb.Build()
}

func buildPullRequestReviewEvent(e *pullRequestReviewEvent) (*bytes.Buffer, error) {
	pb := payloadBuilder{}
	pb.InsertField("アカウント", e.GetSender().GetLogin(), true)
	switch e.GetAction() {
	case "submitted":
		pb.InsertField("アクション", "PR のレビューがされました", true)
		pb.InsertField("タイトル", e.GetPullRequest().GetTitle())
		pb.InsertField("内容", e.GetReview().GetBody())
		pb.InsertField("リンク", e.GetReview().GetHTMLURL())
	default:
		pb.InsertField("アクション", fmt.Sprintf("PullRequestReviewEvent (%s)", e.GetAction()))
	}
	return pb.Build()
}

func buildPullRequestReviewCommentEvent(e *pullRequestReviewCommentEvent) (*bytes.Buffer, error) {
	pb := payloadBuilder{}
	pb.InsertField("アカウント", e.GetSender().GetLogin(), true)
	switch e.GetAction() {
	case "created":
		pb.InsertField("アクション", "PR にコメントされました", true)
		pb.InsertField("コメント", e.GetComment().GetBody())
		pb.InsertField("リンク", e.GetComment().GetHTMLURL())
	case "edited":
		pb.InsertField("アクション", "PR のコメントが変更されました", true)
		if body := e.GetChanges().GetBody(); body != nil {
			pb.InsertField("コメント(変更前)", e.GetChanges().GetBody().GetFrom())
			pb.InsertField("コメント(変更後)", e.GetComment().GetBody())
		}
		pb.InsertField("リンク", e.GetComment().GetHTMLURL())
	case "deleted":
		pb.InsertField("アクション", "PR のコメントが削除されました", true)
		pb.InsertField("コメント", e.GetComment().GetBody())
		pb.InsertField("リンク", e.GetComment().GetHTMLURL())
	default:
		pb.InsertField("アクション", fmt.Sprintf("PullRequestReviewCommentEvent (%s)", e.GetAction()))
	}
	return pb.Build()
}

func buildPullRequestTargetEvent(e *pullRequestTargetEvent) (*bytes.Buffer, error) {
	pb := payloadBuilder{}
	pb.InsertField("アカウント", e.GetSender().GetLogin(), true)
	switch e.GetAction() {
	case "opened":
		pb.InsertField("アクション", "PR がオープンされました", true)
		pb.InsertField("タイトル", e.GetPullRequest().GetTitle())
		pb.InsertField("内容", e.GetPullRequest().GetBody())
		pb.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "edited":
		pb.InsertField("アクション", "PR が編集されました", true)
		if title := e.GetChanges().GetTitle(); title != nil {
			pb.InsertField("タイトル(変更前)", title.GetFrom())
			pb.InsertField("タイトル(変更後)", e.GetPullRequest().GetTitle())
		}
		if body := e.GetChanges().GetBody(); body != nil {
			pb.InsertField("内容(変更前)", body.GetFrom())
			pb.InsertField("内容(変更後)", e.GetPullRequest().GetBody())
		}
		pb.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "closed":
		pb.InsertField("アクション", "PR がクローズされました", true)
		pb.InsertField("タイトル", e.GetPullRequest().GetTitle())
		pb.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "reopened":
		pb.InsertField("アクション", "PR が再オープンされました", true)
		pb.InsertField("タイトル", e.GetPullRequest().GetTitle())
		pb.InsertField("内容", e.GetPullRequest().GetBody())
		pb.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "assigned":
		pb.InsertField("アクション", "PR にアサインされました", true)
		pb.InsertField("対象者", e.GetAssignee().GetLogin())
		pb.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "unassigned":
		pb.InsertField("アクション", "PR にアンアサインされました", true)
		pb.InsertField("対象者", e.GetAssignee().GetLogin())
		pb.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "review_requested":
		pb.InsertField("アクション", "PR のレビューをお願いされました", true)
		pb.InsertField("対象者", e.GetRequestedReviewer().GetLogin())
		pb.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "review_request_removed":
		pb.InsertField("アクション", "PR のレビュー要求が取下げされました", true)
		pb.InsertField("対象者", e.GetRequestedReviewer().GetLogin())
		pb.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "ready_for_review":
		pb.InsertField("アクション", "PR の準備が整いました", true)
		pb.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "labeled":
		pb.InsertField("アクション", "PR にラベルが付与されました", true)
		pb.InsertField("ラベル", e.GetLabel().GetName())
		pb.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "unlabeled":
		pb.InsertField("アクション", "PR のラベルが外されました", true)
		pb.InsertField("ラベル", e.GetLabel().GetName())
		pb.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "locked":
		pb.InsertField("アクション", "PR がロックされました", true)
		pb.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "unlocked":
		pb.InsertField("アクション", "PR のロックが解除されました", true)
		pb.InsertField("リンク", e.GetPullRequest().GetHTMLURL())
	case "synchronize":
		pb.InsertField("アクション", "PullRequestEvent synchronize")
	default:
		pb.InsertField("アクション", fmt.Sprintf("PullRequestTargetEvent (%s)", e.GetAction()))
	}
	return pb.Build()
}

func buildPushEvent(e *pushEvent) (*bytes.Buffer, error) {
	pb := payloadBuilder{}
	pb.InsertField("アカウント", e.GetSender().GetLogin(), true)
	pb.InsertField("アクション", "プッシュされました", true)
	pb.InsertField("対象", e.GetRef())
	if e.Commits != nil {
		var b strings.Builder
		for _, c := range e.Commits {
			b.WriteString(
				fmt.Sprintf("<%s|%s> %s\n", c.GetURL(), c.GetID()[:7], c.GetMessage()),
			)
		}
		pb.InsertField("Commit", b.String())
	}
	pb.InsertField("リンク", e.GetRepo().GetHTMLURL())
	return pb.Build()
}
