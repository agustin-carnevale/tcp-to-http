package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/agustin-carnevale/tcp-to-http/internal/request"
	"github.com/agustin-carnevale/tcp-to-http/internal/response"
)

type Server struct {
	Listener net.Listener
	handler  Handler
	closed   atomic.Bool
}

func Serve(port int, handler Handler) (*Server, error) {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	server := Server{
		Listener: listener,
		handler:  handler,
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
	defer conn.Close()

	// Request
	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Fatalf("Error getting/parsing request: %v", err)
		return
	}

	// Response
	respWriter := &response.Writer{
		Connection: conn,
	}

	s.handler(respWriter, req)
}
