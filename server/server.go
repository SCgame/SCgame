package server

import (
	"io"
	"log"
	"net"
)

type Server struct {
	listener net.Listener
}

func New(addr string) *Server {
	l, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatal(err)
	}

	return &Server{listener: l}
}

func (s *Server) Close() {
	s.listener.Close()
}

func (s *Server) Handle() {
	for {
		conn, err := s.listener.Accept()

		if err != nil {
			log.Fatal(err)
		}

		go func(c net.Conn) {
			io.Copy(c, c)
		}(conn)
	}
}
