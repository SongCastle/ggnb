package message

import "bytes"

type Message interface {
	Init(payload interface{}) error
	NeedToDeliver() bool
	ToPayload() *bytes.Buffer

	ToDummyPayload() *bytes.Buffer
}
