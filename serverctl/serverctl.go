package main

import (
	"scgame/server"
)

func main() {
	s := server.New(":2000")
	defer s.Close()

	s.Handle()
}
