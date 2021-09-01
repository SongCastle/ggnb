package builder

import (
	"bytes"
	"encoding/json"
	"fmt"
)

const (
	Color = "#2eb67d"
	ErrorColor = "#e01e5a"
	Fallback = "GitHub Notifitation"
	TitileLink = "https://github.com/SongCastle/ggnb"
	Title = "GitHub Notification"
)

type field struct {
	Title *string `json:"title"`
	Value *string `json:"value"`
	Short bool `json:"short"`
}

type attachment struct {
	Color *string `json:"color"`
	Fallback *string `json:"fallback"`
	Fields []*field `json:"fields"`
	TitileLink *string `json:"title_link"`
	Title *string `json:"title"`
}

func (a *attachment) InsertField(title, value string, short ...bool) {
	if title != "" && value != "" {
		a.Fields = append(
			a.Fields,
			&field{Title: &title, Value: &value, Short: getShort(short...)},
		)
	}
}

func (a *attachment) Build() (*bytes.Buffer, error) {
	p := payloadBuilder{Attachments: []*attachment{a}}
	return p.Build()
}

type payloadBuilder struct {
	Attachments []*attachment `json:"attachments"`
}

func (pb *payloadBuilder) Build() (*bytes.Buffer, error) {
	j, err := json.Marshal(*pb)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(j), nil
}

func NewAttachment() *attachment {
	return &attachment{
		Color: toP(Color),
		Fallback: toP(Fallback),
		TitileLink: toP(TitileLink),
		Title: toP(Title),
	}
}

func BuildError(err error) (*bytes.Buffer, error) {
	a := NewAttachment()
	a.Color = toP(ErrorColor)
	a.InsertField("エラー", fmt.Sprintf("%v", err))
	return a.Build()
}

func getShort(short ...bool) bool {
	if len(short) == 0 {
		return false
	}
	return short[0]
}

func toP(s string) *string {
	return &s
}
