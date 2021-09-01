package main

import "github.com/SongCastle/ggnb/handler"

func main() {
	h := handler.New()
	h.Init()
	h.Start()
}
