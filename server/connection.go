package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

type Connection struct {
	server  *Server
	game    *Game
	tcpConn *net.TCPConn
	name    string
}

func NewConnection(s *Server, g *Game, c *net.TCPConn) *Connection {
	return &Connection{server: s, game: g, tcpConn: c}
}

func (c *Connection) Handle() {
	login, err := c.read()
	if err != nil {
		c.close("ERROR 0")
		return
	}

	succ := c.authorize(login)
	if !succ {
		c.close("ERROR 1 not authenticated")
		return
	}

	for {
		message, err := c.read()
		if err != nil {
			c.close("ERROR 0")
			return
		}
		c.log(message)

		req := NewRequest(c.name, message)
		c.game.RequestChan <- req

		res := (<-req.ResponseChan).Content
		close(req.ResponseChan)

		c.ok(res)
	}
}

func (c *Connection) read() (message string, err error) {
	c.tcpConn.SetReadDeadline(time.Now().Add(timeout))
	message, err = bufio.NewReader(c.tcpConn).ReadString('\n')
	return
}

func (c *Connection) ok(message string) {
	fmt.Printf("%p: OK %s\n", c, message)
	fmt.Fprintf(c.tcpConn, "OK %s\n", message)
}

func (c *Connection) close(errMessage string) {
	fmt.Printf("%p: -> %s\n", c, errMessage)
	fmt.Fprintf(c.tcpConn, "%s\n", errMessage)
	c.tcpConn.Close()
	// notify the server
}

func (c *Connection) log(message string) {
	fmt.Printf("%p: <- %s\n", c, message)
}

func (c *Connection) authorize(command string) bool {
	tokens := strings.Split(command, " ")
	if tokens[0] != "LOGIN" || len(tokens) != 2 {
		// notify the server
		return false
	}

	// notify the server
	c.name = tokens[1]
	fmt.Printf("%p: Authorized %s\n", c, c.name)
	fmt.Fprintf(c.tcpConn, "OK\n")
	return true
}
