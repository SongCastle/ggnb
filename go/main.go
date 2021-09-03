package main

import (
	"fmt"

	"github.com/SongCastle/ggnb/income/message"
	"github.com/SongCastle/ggnb/outcome/client"
	"github.com/SongCastle/ggnb/handler"
)

func main() {
	// Create Message (income)
	m, err := message.NewMessage()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	// Create Client (outcome)
	c := client.NewClient()
	if err := c.Init(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	// Create & Start Handler
	h := handler.New(m, c)
	h.Start()
}
