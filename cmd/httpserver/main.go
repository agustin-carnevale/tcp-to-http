package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/agustin-carnevale/tcp-to-http/internal/request"
	"github.com/agustin-carnevale/tcp-to-http/internal/response"
	"github.com/agustin-carnevale/tcp-to-http/internal/server"
)

const port = 42069

func handler(w io.Writer, req *request.Request) *server.HandlerError {

	// req.Print()
	if req.RequestLine.RequestTarget == "/yourproblem" {
		return &server.HandlerError{
			StatusCode: response.StatusBadRequest,
			Message:    "Your problem is not my problem\n",
		}
	}

	if req.RequestLine.RequestTarget == "/myproblem" {
		return &server.HandlerError{
			StatusCode: response.StatusInternalServerError,
			Message:    "Woopsie, my bad\n",
		}
	}

	_, err := w.Write([]byte("All good, frfr\n"))
	if err != nil {
		log.Fatalf("Error writing response body to buffer: %v", err)
	}

	return nil
}

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
