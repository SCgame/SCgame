package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

const timeout time.Duration = 5 * time.Second

type Server struct {
	listener *net.TCPListener
	game     *Game
}

func New(addr string) *Server {
	l, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatal(err)
	}

	return &Server{listener: l.(*net.TCPListener), game: NewGame()}
}

func (s *Server) Close() {
	s.listener.Close()
}

func (s *Server) Handle() {
	for {
		var conn *net.TCPConn
		conn, err := s.listener.AcceptTCP()

		if err != nil {
			log.Fatal(err)
		}

		go func(c *net.TCPConn) {
			for {
				c.SetReadDeadline(time.Now().Add(timeout))

				message, err := bufio.NewReader(c).ReadString('\n')
				if err != nil {
					fmt.Println("Closing connection")
					c.Close()
					return
				}

				req := NewRequest(message)
				s.game.RequestChan <- req

				response := (<-req.ResponseChan).Content
				fmt.Println(response)
				fmt.Fprintf(c, response)
			}
		}(conn)
	}
}
