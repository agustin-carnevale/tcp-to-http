package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/agustin-carnevale/tcp-to-http/internal/request"
	"github.com/agustin-carnevale/tcp-to-http/internal/response"
	"github.com/agustin-carnevale/tcp-to-http/internal/server"
)

const port = 42069

func handler(w *response.Writer, req *request.Request) {
	var statusCode response.StatusCode = response.StatusOK
	var html string = ""

	if req.RequestLine.RequestTarget == "/yourproblem" {
		statusCode = response.StatusBadRequest
		html = `<html>
		<head>
			<title>400 Bad Request</title>
		</head>
		<body>
			<h1>Bad Request</h1>
			<p>Your request honestly kinda sucked.</p>
		</body>
		</html>`
	} else if req.RequestLine.RequestTarget == "/myproblem" {
		statusCode = response.StatusInternalServerError
		html = `<html>
		<head>
			<title>500 Internal Server Error</title>
		</head>
		<body>
			<h1>Internal Server Error</h1>
			<p>Okay, you know what? This one is on me.</p>
		</body>
		</html>`
	} else {
		html = `<html>
		<head>
			<title>200 OK</title>
		</head>
		<body>
			<h1>Success!</h1>
			<p>Your request was an absolute banger.</p>
		</body>
		</html>`
	}

	err := w.WriteStatusLine(statusCode)
	if err != nil {
		log.Fatalf("Error writing response status-line: %v", err)
		return
	}

	headers := response.GetDefaultHeaders(len(html))
	headers.SetWithOverride("Content-Type", "text/html")

	err = w.WriteHeaders(headers)
	if err != nil {
		log.Fatalf("Error writing response headers: %v", err)
		return
	}

	_, err = w.WriteBody([]byte(html))
	if err != nil {
		log.Fatalf("Error writing response body: %v", err)
		return
	}

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
