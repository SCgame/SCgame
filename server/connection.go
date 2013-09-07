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
	Name    string
}

func NewConnection(s *Server, g *Game, c *net.TCPConn) *Connection {
	s.ConnMutex.Lock()
	s.Connections += 1
	s.ConnMutex.Unlock()
	return &Connection{server: s, game: g, tcpConn: c}
}

func (c *Connection) Handle() {
	login, err := c.read()
	if err != nil {
		c.Close("ERROR 0")
		return
	}

	succ := c.authorize(login)
	if !succ {
		c.Close("ERROR 1 not authenticated")
		return
	}

	for {
		message, err := c.read()
		if err != nil {
			c.Close("ERROR 0")
			return
		}
		c.log(message)

		req := NewRequest(c.Name, message)
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

func (c *Connection) Close(errMessage string) {
	fmt.Printf("%p: -> %s\n", c, errMessage)
	fmt.Fprintf(c.tcpConn, "%s\n", errMessage)
	c.tcpConn.Close()
	c.server.ConnMutex.Lock()
	c.server.Connections -= 1
	if c.Name != "" {
		delete(c.server.ConnectionPool, c.Name)
	}
	c.server.ConnMutex.Unlock()
}

func (c *Connection) log(message string) {
	fmt.Printf("%p: <- %s\n", c, message)
}

func (c *Connection) authorize(command string) bool {
	tokens := strings.Split(command, " ")
	if tokens[0] != "LOGIN" || len(tokens) != 2 {
		return false
	}

	c.Name = tokens[1]

	c.server.ConnMutex.Lock()
	oldConn, oldPresent := c.server.ConnectionPool[c.Name]
	c.server.ConnectionPool[c.Name] = c
	c.server.ConnMutex.Unlock()
	if oldPresent {
		oldConn.Close("ERROR 2 connection reset by peer")
	}

	fmt.Printf("%p: Authorized %s\n", c, c.Name)
	fmt.Fprintf(c.tcpConn, "OK\n")
	return true
}
