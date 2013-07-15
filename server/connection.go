package server

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

type Connection struct {
	game    *Game
	tcpConn *net.TCPConn
}

func NewConnection(g *Game, c *net.TCPConn) *Connection {
	return &Connection{game: g, tcpConn: c}
}

func (c *Connection) Handle() {
	for {
		c.tcpConn.SetReadDeadline(time.Now().Add(timeout))

		message, err := bufio.NewReader(c.tcpConn).ReadString('\n')
		if err != nil {
			c.tcpConn.Close()
			return
		}

		req := NewRequest(message)
		c.game.RequestChan <- req

		res := (<-req.ResponseChan).Content
		close(req.ResponseChan)

		fmt.Println(res)
		fmt.Fprintf(c.tcpConn, res)
	}

}
