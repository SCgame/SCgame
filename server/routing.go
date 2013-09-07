package server

type Response struct {
	Content string
}

type Request struct {
	Command      string
	User         string
	ResponseChan chan *Response
}

type Game struct {
	RequestChan   chan *Request
	BoardInstance *Board
}

func NewRequest(user string, cmd string) *Request {
	return &Request{User: user, Command: cmd, ResponseChan: make(chan *Response)}
}

func NewGame() *Game {
	board := NewBoard()
	game := &Game{RequestChan: make(chan *Request), BoardInstance: board}
	go game.receiveCommands()
	return game
}

func (g *Game) receiveCommands() {
	for req := range g.RequestChan {
		req.ResponseChan <- g.BoardInstance.RunCommand(req)
	}
}
