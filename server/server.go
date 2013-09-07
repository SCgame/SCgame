package server

import (
	"log"
	"net"
	"sync"
	"time"
)

const timeout time.Duration = 5 * time.Second
const rate time.Duration = 1 * time.Second
const maxConns int = 10

type Server struct {
	listener       *net.TCPListener
	game           *Game
	ConnMutex      sync.Mutex
	Connections    int
	ConnectionPool map[string]*Connection
}

func New(addr string) *Server {
	l, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatal(err)
	}

	return &Server{listener: l.(*net.TCPListener), game: NewGame(),
		Connections:    0,
		ConnectionPool: make(map[string]*Connection)}
}

func (s *Server) Close() {
	s.listener.Close()
}

func (s *Server) Handle() {
	for {
		c, err := s.listener.AcceptTCP()
		if s.Connections >= maxConns {
			c.Close()
			continue
		}

		if err != nil {
			log.Fatal(err)
		}

		conn := NewConnection(s, s.game, c)
		go conn.Handle()
	}
}
