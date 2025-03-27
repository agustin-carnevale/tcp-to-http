package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/agustin-carnevale/tcp-to-http/internal/request"
	"github.com/agustin-carnevale/tcp-to-http/internal/response"
	"github.com/agustin-carnevale/tcp-to-http/internal/server"
)

const port = 42069

func handler(w *response.Writer, req *request.Request) {
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
		handlerProxy(w, req)
		return
	} else if req.RequestLine.RequestTarget == "/yourproblem" {
		handlerYourProblem(w, req)
		return
	} else if req.RequestLine.RequestTarget == "/myproblem" {
		handlerMyProblem(w, req)
		return
	} else if req.RequestLine.RequestTarget == "/video" {
		handlerGetVideoStream(w, req)
		return
	}

	handlerStatusOk(w, req)
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
