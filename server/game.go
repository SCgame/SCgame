package server

import (
	"fmt"
	"strings"
)

type command struct {
	Arity  int
	Method func([]string, *Player) string
}

type Board struct {
	Commands map[string]*command
	Players  map[string]*Player
}

type Position struct {
	X, Y int
}

type Player struct {
	Pos Position
}

func NewPlayer() *Player {
	return &Player{}
}

func (g *Board) getPlayer(req *Request) *Player {
	_, ok := g.Players[req.User]
	if !ok {
		g.Players[req.User] = NewPlayer()
	}
	return g.Players[req.User]
}

func NewBoard() *Board {
	board := &Board{Commands: map[string]*command{}, Players: map[string]*Player{}}
	board.registerCommand("MOVE", 1, func(args []string, player *Player) string {
		direction := args[0]
		switch direction {
		case "UP":
			player.Pos.Y -= 1
		case "DOWN":
			player.Pos.Y += 1
		case "LEFT":
			player.Pos.X -= 1
		case "RIGHT":
			player.Pos.X += 1
		default:
			return "ERROR 103 INVALID DIRECTION"
		}
		return fmt.Sprintf("OK %d %d", player.Pos.X, player.Pos.Y)
	})
	return board
}

func (g *Board) RunCommand(req *Request) *Response {
	cmdWithArgs := strings.Split(req.Command, " ")
	method, ok := g.Commands[cmdWithArgs[0]]
	fmt.Println(cmdWithArgs)

	if !ok {
		return &Response{"ERROR 101 INVALID COMMAND"}
	}
	if len(cmdWithArgs[1:]) != method.Arity {
		return &Response{"ERROR 102 INVALID NUMBER OF ARGUMENTS"}
	}
	return &Response{method.Method(cmdWithArgs[1:], g.getPlayer(req))}
}

func (g *Board) registerCommand(cmd string, arity int, method func([]string, *Player) string) {
	g.Commands[cmd] = &command{arity, method}
}
