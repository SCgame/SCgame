package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type command struct {
	Arity  int
	Method func([]string) string
}

type game struct {
	Commands map[string]command
}

func BuildGame() game {
	return game{map[string]command{}}
}

func main() {
	g := new(game)
	for {
		cmd, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		fmt.Println(g.runCommand(cmd))
	}
}

func (g *game) runCommand(cmd string) string {
	cmdWithArgs := strings.Split(cmd, " ")
	method, ok := g.Commands[cmdWithArgs[0]]
	if !ok {
		return "ERROR 1 INVALID COMMAND"
	}
	if len(cmdWithArgs[1:]) != method.Arity {
		return "ERROR 2 INVALID NUMBER OF ARGUMENTS"
	}
	return method.Method(cmdWithArgs[1:])
}

func (g *game) registerCommand(cmd string, arity int, method func([]string) string) {
	g.Commands[cmd] = command{arity, method}
}
