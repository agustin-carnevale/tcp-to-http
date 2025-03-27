package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/agustin-carnevale/tcp-to-http/internal/headers"
	"github.com/agustin-carnevale/tcp-to-http/internal/request"
	"github.com/agustin-carnevale/tcp-to-http/internal/response"
	"github.com/agustin-carnevale/tcp-to-http/internal/server"
)

const port = 42069

// const endOfChunkedBody = "0\r\n\r\n"

func handler(w *response.Writer, req *request.Request) {
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
		proxyHandler(w, req)
		return
	}

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

func proxyHandler(w *response.Writer, req *request.Request) {
	proxyToTarget := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin")
	proxyToUrl := "https://httpbin.org" + proxyToTarget

	fmt.Println(proxyToUrl)

	resp, err := http.Get(proxyToUrl)
	if err != nil {
		// http.Error(w, "Failed to reach target", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	headers := headers.Headers{
		"connection":        "close",
		"content-type":      "text/plain",
		"transfer-encoding": "chunked",
	}

	w.WriteStatusLine(response.StatusOK)
	w.WriteHeaders(headers)

	buf := make([]byte, 1024)

	for {
		n, err := resp.Body.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			// http.Error(w, "Error reading from target", http.StatusInternalServerError)
			return
		}
		// fmt.Println("Bytes read:", n)

		// if strings.Contains(string(buf[:n]), endOfChunkedBody) {
		// 	fmt.Println("END OF BODY RECEIVED!!")
		// }

		// Write chunked response
		w.WriteChunkedBody(buf[:n])

	}
	// fmt.Println("End of Body")
	w.WriteChunkedBodyDone()

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
