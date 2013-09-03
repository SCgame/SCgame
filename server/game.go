package server

import "fmt"

type Response struct {
	Content string
}

type Request struct {
	Command      string
	User         string
	ResponseChan chan *Response
}

type Game struct {
	RequestChan chan *Request
}

func NewRequest(user string, cmd string) *Request {
	return &Request{User: user, Command: cmd, ResponseChan: make(chan *Response)}
}

func NewGame() *Game {
	game := &Game{make(chan *Request)}
	go game.receiveCommands()
	return game
}

func (g *Game) receiveCommands() {
	for req := range g.RequestChan {
		req.ResponseChan <- &Response{fmt.Sprintf("RESPONSE: %s", req.Command)}
	}
}
