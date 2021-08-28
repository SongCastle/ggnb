package outcome

import (
	"bytes"
	"fmt"
)

func new() client {
	sc := slackClient{}
	sc.Init()
	return &sc
}

func Send(msg *bytes.Buffer) error {
	if msg == nil {
		fmt.Println("Skipped\n")
		return nil
	}
	fmt.Printf("payload: %s\n", msg.String())

	client := new()
	if err := client.Post(msg); err != nil {
		return err
	}
	return nil
}
