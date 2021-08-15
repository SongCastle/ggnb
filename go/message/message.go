package message

import "bytes"

type Message interface {
	Init(payload interface{}) error
	ToPayload() (*bytes.Buffer, error)
	ToDummyPayload() *bytes.Buffer
}
