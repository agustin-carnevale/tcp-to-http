package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/agustin-carnevale/tcp-to-http/internal/response"
)

type Server struct {
	Listener net.Listener
	closed   atomic.Bool
}

func Serve(port int) (*Server, error) {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	server := Server{
		Listener: listener,
	}

	go server.listen()

	return &server, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.Listener != nil {
		return s.Listener.Close()
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			fmt.Println("Error while accepting tcp connection:", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	headers := response.GetDefaultHeaders(0)

	err := response.WriteStatusLine(conn, response.StatusOK)
	if err != nil {
		log.Fatalf("Error writing response status-line: %v", err)
	}

	response.WriteHeaders(conn, headers)
	if err != nil {
		log.Fatalf("Error writing response headers: %v", err)
	}

	err = conn.Close()
	if err != nil {
		log.Fatalf("Error closing connection: %v", err)
	}
}
