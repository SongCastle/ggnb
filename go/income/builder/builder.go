package builder

import (
	"bytes"
	"encoding/json"
	"fmt"
)

const (
	Fallback = "GitHub Notifitation"
	Color = "#2eb67d"
	Title = "GitHub Notification"
	TitileLink = "https://github.com/SongCastle/ggnb"
	ErrorColor = "#e01e5a"
)

type field struct {
	Title *string `json:"title"`
	Value *string `json:"value"`
	Short bool `json:"short"`
}

type attachment struct {
	Fallback *string `json:"fallback"`
	Color *string `json:"color"`
	Title *string `json:"title"`
	TitileLink *string `json:"title_link"`
	Fields []*field `json:"fields"`
}

func (a *attachment) InsertField(title, value string, short ...bool) {
	a.Fields = append(
		a.Fields,
		&field{Title: &title, Value: &value, Short: getShort(short...)},
	)
}

func (a *attachment) Build() (*bytes.Buffer, error) {
	p := payloadBuilder{Attachment: []*attachment{a}}
	return p.Build()
}

type payloadBuilder struct {
	Attachment []*attachment `json:"attachments"`
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
		Fallback: sToP(Fallback),
		Color: sToP(Color),
		Title: sToP(Title),
		TitileLink: sToP(TitileLink),
	}
}

func BuildError(err error) (*bytes.Buffer, error) {
	a := NewAttachment()
	a.Color = sToP(ErrorColor)
	a.InsertField("エラー", fmt.Sprintf("%v", err))
	return a.Build()
}

func getShort(short ...bool) bool {
	if len(short) == 0 {
		return false
	}
	return short[0]
}

func sToP(s string) *string {
	return &s
}
