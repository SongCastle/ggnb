package builder

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAttachmentInsertField(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	a := NewAttachment()
	a.InsertField("Account", "Codertocat")
	a.InsertField("Action", "Push", true)
	a.InsertField("Link", "https://github.com/SongCastle/ggnb", false)
	a.InsertField("", "")

	assert.Equal(len(a.Fields), 3)
	assert.Equal(a.Fields[0].Title, toP("Account"))
	assert.Equal(a.Fields[0].Value, toP("Codertocat"))
	assert.Equal(a.Fields[0].Short, false)
	assert.Equal(a.Fields[1].Title, toP("Action"))
	assert.Equal(a.Fields[1].Value, toP("Push"))
	assert.Equal(a.Fields[1].Short, true)
	assert.Equal(a.Fields[2].Title, toP("Link"))
	assert.Equal(a.Fields[2].Value, toP("https://github.com/SongCastle/ggnb"))
	assert.Equal(a.Fields[2].Short, false)
}

func TestAttachmentBuild(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	a := NewAttachment()
	a.InsertField("Account", "Codertocat")
	a.InsertField("Action", "Push", true)
	a.InsertField("Link", "https://github.com/SongCastle/ggnb", false)

	buf, err := a.Build()
	assert.Nil(err)

	msg :=
		fmt.Sprintf(
			`
			{
				"attachments":
					[
						{
							"color":"%s",
							"fallback":"%s",
							"fields":[
								{"title":"Account","value":"Codertocat","short":false},
								{"title":"Action","value":"Push","short":true},
								{"title":"Link","value":"https://github.com/SongCastle/ggnb","short":false}
							],
							"title_link":"%s",
							"title":"%s"
						}
					]
			}
		`, Color, Fallback, TitileLink, Title,
		)
	msg = strings.ReplaceAll(msg, "\t", "")
	msg = strings.ReplaceAll(msg, "\n", "")

	assert.Equal(buf, bytes.NewBufferString (msg))
}

func TestPayloadBuilderBuild(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	a := NewAttachment()
	a.InsertField("Account", "Codertocat")
	a.InsertField("Action", "Push", true)
	a.InsertField("Link", "https://github.com/SongCastle/ggnb", false)

	pb := &payloadBuilder{
		Attachments: []*attachment{
			a,
		},
	}

	buf, err := pb.Build()
	assert.Nil(err)

	msg :=
		fmt.Sprintf(
			`
				{
					"attachments":
						[
							{
								"color":"%s",
								"fallback":"%s",
								"fields":[
									{"title":"Account","value":"Codertocat","short":false},
									{"title":"Action","value":"Push","short":true},
									{"title":"Link","value":"https://github.com/SongCastle/ggnb","short":false}
								],
								"title_link":"%s",
								"title":"%s"
							}
						]
				}
			`, Color, Fallback, TitileLink, Title,
		)
	msg = strings.ReplaceAll(msg, "\t", "")
	msg = strings.ReplaceAll(msg, "\n", "")

	assert.Equal(buf, bytes.NewBufferString (msg))
}

func TestNewAttachment(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	a := NewAttachment()
	assert.IsType(a, &attachment{})
	assert.Equal(a.Color, toP(Color))
	assert.Equal(a.Fallback, toP(Fallback))
	assert.Nil(a.Fields)
	assert.Equal(a.TitileLink, toP(TitileLink))
	assert.Equal(a.Title, toP(Title))
}

func TestBuildError(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	arg := errors.New("test")
	buf, err := BuildError(arg)
	assert.Nil(err)

	msg :=
		fmt.Sprintf(
			`
				{
					"attachments":
						[
							{
								"color":"%s",
								"fallback":"%s",
								"fields":[
									{"title":"エラー","value":"%v","short":false}
								],
								"title_link":"%s",
								"title":"%s"
							}
						]
				}
			`, ErrorColor, Fallback, arg, TitileLink, Title,
		)
	msg = strings.ReplaceAll(msg, "\t", "")
	msg = strings.ReplaceAll(msg, "\n", "")

	assert.Equal(buf, bytes.NewBufferString (msg))
}

func TestGetShort(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	t.Run("without args", func(t *testing.T){
		assert.False(getShort())
	})

	t.Run("with an arg", func(t *testing.T){
		assert.True(getShort(true))
		assert.False(getShort(false))
	})

	t.Run("with some args", func(t *testing.T){
		assert.True(getShort(true, true))
		assert.True(getShort(true, false))
		assert.False(getShort(false, true))
		assert.False(getShort(false, false))
	})
}

func TestToP(t *testing.T) {
	t.Parallel()

	s := "test"
	assert.Equal(t, toP(s), &s)
}
