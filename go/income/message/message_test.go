package message

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"unsafe"

	"github.com/SongCastle/ggnb/income/builder"
	"github.com/stretchr/testify/assert"
)

func TestGitHubMessageInit(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	t.Run("with nil", func(t *testing.T) {
		gm := GitHubMessage{}
		err := gm.Init(nil)
		assert.EqualError(err, "invalid payload")
	})

	t.Run("without headers", func(t *testing.T) {
		gm := GitHubMessage{}
		err := gm.Init(fields{})
		assert.EqualError(err, "missing headers")
	})

	t.Run("without GitHub Event header", func(t *testing.T) {
		gm := GitHubMessage{}
		err := gm.Init(fields{"headers": fields{}})
		assert.EqualError(err, fmt.Sprintf("missing %s header", EventHeader))
	})

	t.Run("with unknown GitHub Event header", func(t *testing.T) {
		gm := GitHubMessage{}
		err := gm.Init(fields{"headers": fields{EventHeader: "xxxxx"}, "body": "{}"})
		assert.EqualError(err, "unknown X-Github-Event in message: xxxxx")
	})

	t.Run("without body", func(t *testing.T) {
		gm := GitHubMessage{}
		err := gm.Init(fields{"headers": fields{EventHeader: "push"}})
		assert.EqualError(err, "missing body")
	})

	t.Run("with body", func(t *testing.T) {
		gm := GitHubMessage{}
		err := gm.Init(fields{"headers": fields{EventHeader: "push"}, "body": "{}"})
		assert.Nil(err)
		_, ok := gm.event.(*pushEvent)
		assert.True(ok)
	})
}

func TestGitHubMessageToPayload(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	t.Run("commit_comment", func(t *testing.T) {
		gm := GitHubMessage{}
		json, err := os.ReadFile("./test_data/commit_comment.json")
		if err != nil {
			t.Error(err)
		}

		err = gm.Init(
			fields{
				"headers": fields{EventHeader: "commit_comment"},
				"body": *(*string)(unsafe.Pointer(&json)),
			},
		)
		assert.Nil(err)

		buf, err := gm.ToPayload()
		assert.Nil(err)
		assert.IsType(buf, &bytes.Buffer{})

		a := builder.NewAttachment()
		a.InsertField("アカウント", "Codertocat", true)
		a.InsertField("アクション", "コメントされました", true)
		a.InsertField("コメント", "This is a really good change! :+1:")
		a.InsertField("CommitID", "6113728f27ae82c7b1a177c8d03f9e96e0adf246")
		a.InsertField("リンク", "https://github.com/Codertocat/Hello-World/commit/6113728f27ae82c7b1a177c8d03f9e96e0adf246#commitcomment-33548674")
		ebuf, err := a.Build()
		if err != nil {
			t.Error(err)
		}

		assert.Equal(buf, ebuf)
	})

	t.Run("create", func(t *testing.T) {
		gm := GitHubMessage{}
		json, err := os.ReadFile("./test_data/create.json")
		if err != nil {
			t.Error(err)
		}

		err = gm.Init(
			fields{
				"headers": fields{EventHeader: "create"},
				"body": *(*string)(unsafe.Pointer(&json)),
			},
		)
		assert.Nil(err)

		buf, err := gm.ToPayload()
		assert.Nil(err)
		assert.IsType(buf, &bytes.Buffer{})

		a := builder.NewAttachment()
		a.InsertField("アカウント", "Codertocat", true)
		a.InsertField("アクション", "タグが作成されました", true)
		a.InsertField("タグ名", "simple-tag")
		a.InsertField("リンク", "https://github.com/Codertocat/Hello-World")
		ebuf, err := a.Build()
		if err != nil {
			t.Error(err)
		}

		assert.Equal(buf, ebuf)
	})

	t.Run("delete", func(t *testing.T) {
		gm := GitHubMessage{}
		json, err := os.ReadFile("./test_data/delete.json")
		if err != nil {
			t.Error(err)
		}

		err = gm.Init(
			fields{
				"headers": fields{EventHeader: "delete"},
				"body": *(*string)(unsafe.Pointer(&json)),
			},
		)
		assert.Nil(err)

		buf, err := gm.ToPayload()
		assert.Nil(err)
		assert.IsType(buf, &bytes.Buffer{})

		a := builder.NewAttachment()
		a.InsertField("アカウント", "Codertocat", true)
		a.InsertField("アクション", "タグが削除されました", true)
		a.InsertField("タグ名", "simple-tag")
		a.InsertField("リンク", "https://github.com/Codertocat/Hello-World")
		ebuf, err := a.Build()
		if err != nil {
			t.Error(err)
		}

		assert.Equal(buf, ebuf)
	})

	t.Run("issue_comment", func(t *testing.T) {
		gm := GitHubMessage{}
		json, err := os.ReadFile("./test_data/issue_comment.json")
		if err != nil {
			t.Error(err)
		}

		err = gm.Init(
			fields{
				"headers": fields{EventHeader: "issue_comment"},
				"body": *(*string)(unsafe.Pointer(&json)),
			},
		)
		assert.Nil(err)

		buf, err := gm.ToPayload()
		assert.Nil(err)
		assert.IsType(buf, &bytes.Buffer{})

		a := builder.NewAttachment()
		a.InsertField("アカウント", "Codertocat", true)
		a.InsertField("アクション", "コメントされました", true)
		a.InsertField("コメント", "You are totally right! I'll get this fixed right away.")
		a.InsertField("リンク", "https://github.com/Codertocat/Hello-World/issues/1#issuecomment-492700400")
		ebuf, err := a.Build()
		if err != nil {
			t.Error(err)
		}

		assert.Equal(buf, ebuf)
	})

	t.Run("issues", func(t *testing.T) {
		gm := GitHubMessage{}
		json, err := os.ReadFile("./test_data/issues.json")
		if err != nil {
			t.Error(err)
		}

		err = gm.Init(
			fields{
				"headers": fields{EventHeader: "issues"},
				"body": *(*string)(unsafe.Pointer(&json)),
			},
		)
		assert.Nil(err)

		buf, err := gm.ToPayload()
		assert.Nil(err)
		assert.IsType(buf, &bytes.Buffer{})

		a := builder.NewAttachment()
		a.InsertField("アカウント", "Codertocat", true)
		a.InsertField("アクション", "Issue が編集されました", true)
		a.InsertField("タイトル(変更前)", "Spelling error in the README")
		a.InsertField("タイトル(変更後)", "Spelling error in the README file")
		a.InsertField("リンク", "https://github.com/Codertocat/Hello-World/issues/1")
		ebuf, err := a.Build()
		if err != nil {
			t.Error(err)
		}

		assert.Equal(buf, ebuf)
	})

	t.Run("pull_request", func(t *testing.T) {
		gm := GitHubMessage{}
		json, err := os.ReadFile("./test_data/pull_request.json")
		if err != nil {
			t.Error(err)
		}

		err = gm.Init(
			fields{
				"headers": fields{EventHeader: "pull_request"},
				"body": *(*string)(unsafe.Pointer(&json)),
			},
		)
		assert.Nil(err)

		buf, err := gm.ToPayload()
		assert.Nil(err)
		assert.IsType(buf, &bytes.Buffer{})

		a := builder.NewAttachment()
		a.InsertField("アカウント", "Codertocat", true)
		a.InsertField("アクション", "PR がオープンされました", true)
		a.InsertField("タイトル", "Update the README with new information.")
		a.InsertField("内容", "This is a pretty simple change that we need to pull into master.")
		a.InsertField("リンク", "https://github.com/Codertocat/Hello-World/pull/2")
		ebuf, err := a.Build()
		if err != nil {
			t.Error(err)
		}

		assert.Equal(buf, ebuf)
	})

	t.Run("pull_request_review", func(t *testing.T) {
		gm := GitHubMessage{}
		json, err := os.ReadFile("./test_data/pull_request_review.json")
		if err != nil {
			t.Error(err)
		}

		err = gm.Init(
			fields{
				"headers": fields{EventHeader: "pull_request_review"},
				"body": *(*string)(unsafe.Pointer(&json)),
			},
		)
		assert.Nil(err)

		buf, err := gm.ToPayload()
		assert.Nil(err)
		assert.IsType(buf, &bytes.Buffer{})

		a := builder.NewAttachment()
		a.InsertField("アカウント", "Codertocat", true)
		a.InsertField("アクション", "PR のレビューがされました", true)
		a.InsertField("タイトル", "Update the README with new information.")
		a.InsertField("内容", "LGTM")
		a.InsertField("リンク", "https://github.com/Codertocat/Hello-World/pull/2#pullrequestreview-237895671")
		ebuf, err := a.Build()
		if err != nil {
			t.Error(err)
		}

		assert.Equal(buf, ebuf)
	})

	t.Run("pull_request_review_comment", func(t *testing.T) {
		gm := GitHubMessage{}
		json, err := os.ReadFile("./test_data/pull_request_review_comment.json")
		if err != nil {
			t.Error(err)
		}

		err = gm.Init(
			fields{
				"headers": fields{EventHeader: "pull_request_review_comment"},
				"body": *(*string)(unsafe.Pointer(&json)),
			},
		)
		assert.Nil(err)

		buf, err := gm.ToPayload()
		assert.Nil(err)
		assert.IsType(buf, &bytes.Buffer{})

		a := builder.NewAttachment()
		a.InsertField("アカウント", "Codertocat", true)
		a.InsertField("アクション", "PR にコメントされました", true)
		a.InsertField("コメント", "Maybe you should use more emoji on this line.")
		a.InsertField("リンク", "https://github.com/Codertocat/Hello-World/pull/2#discussion_r284312630")
		ebuf, err := a.Build()
		if err != nil {
			t.Error(err)
		}

		assert.Equal(buf, ebuf)
	})

	t.Run("pull_request_target", func(t *testing.T) {
		gm := GitHubMessage{}
		json, err := os.ReadFile("./test_data/pull_request_target.json")
		if err != nil {
			t.Error(err)
		}

		err = gm.Init(
			fields{
				"headers": fields{EventHeader: "pull_request_target"},
				"body": *(*string)(unsafe.Pointer(&json)),
			},
		)
		assert.Nil(err)

		buf, err := gm.ToPayload()
		assert.Nil(err)
		assert.IsType(buf, &bytes.Buffer{})

		a := builder.NewAttachment()
		a.InsertField("アカウント", "Codertocat", true)
		a.InsertField("アクション", "PR がオープンされました", true)
		a.InsertField("タイトル", "Update the README with new information.")
		a.InsertField("内容", "This is a pretty simple change that we need to pull into master.")
		a.InsertField("リンク", "https://github.com/Codertocat/Hello-World/pull/2")
		ebuf, err := a.Build()
		if err != nil {
			t.Error(err)
		}

		assert.Equal(buf, ebuf)
	})

	t.Run("push", func(t *testing.T) {
		gm := GitHubMessage{}
		json, err := os.ReadFile("./test_data/push.json")
		if err != nil {
			t.Error(err)
		}

		err = gm.Init(
			fields{
				"headers": fields{EventHeader: "push"},
				"body": *(*string)(unsafe.Pointer(&json)),
			},
		)
		assert.Nil(err)

		buf, err := gm.ToPayload()
		assert.Nil(err)
		assert.IsType(buf, &bytes.Buffer{})

		a := builder.NewAttachment()
		a.InsertField("アカウント", "Codertocat", true)
		a.InsertField("アクション", "プッシュされました", true)
		a.InsertField("対象", "refs/tags/simple-tag")
		a.InsertField("Commit",
			"<https://github.com/Codertocat/Hello-World/commit/0123456789012345678901234567890123456789|0123456> Small Changes\n",
		)
		a.InsertField("リンク", "https://github.com/Codertocat/Hello-World")
		ebuf, err := a.Build()
		if err != nil {
			t.Error(err)
		}

		assert.Equal(buf, ebuf)
	})

	t.Run("NOT targeted event", func(t *testing.T) {
		gm := GitHubMessage{}
		json, err := os.ReadFile("./test_data/not_targeted.json")
		if err != nil {
			t.Error(err)
		}

		err = gm.Init(
			fields{
				"headers": fields{EventHeader: "check_run"},
				"body": *(*string)(unsafe.Pointer(&json)),
			},
		)
		assert.Nil(err)

		buf, err := gm.ToPayload()
		assert.Nil(err)
		assert.Nil(buf)
	})
}

func TestGitHubMessageToDummyPayload(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	gm := GitHubMessage{}
	buf, err := gm.ToDummyPayload()
	assert.Nil(err)
	assert.IsType(buf, &bytes.Buffer{})

	a := builder.NewAttachment()
	a.InsertField("アカウント", "bot", true)
	a.InsertField("アクション", "debug", true)
	a.InsertField("内容", "ok")
	ebuf, err := a.Build()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(buf, ebuf)
}
