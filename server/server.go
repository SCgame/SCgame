package server

import (
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
		c, err := s.listener.AcceptTCP()

		if err != nil {
			log.Fatal(err)
		}

		conn := NewConnection(s.game, c)
		go conn.Handle()
	}
}
