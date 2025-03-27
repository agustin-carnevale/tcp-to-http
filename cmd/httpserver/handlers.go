package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/agustin-carnevale/tcp-to-http/internal/headers"
	"github.com/agustin-carnevale/tcp-to-http/internal/request"
	"github.com/agustin-carnevale/tcp-to-http/internal/response"
)

func handlerStatusOk(w *response.Writer, req *request.Request) {
	statusCode := response.StatusOK
	html := `<html>
		<head>
			<title>200 OK</title>
		</head>
		<body>
			<h1>Success!</h1>
			<p>Your request was an absolute banger.</p>
		</body>
		</html>`

	writeHTMLResponse(w, statusCode, html)
}

func handlerMyProblem(w *response.Writer, req *request.Request) {
	statusCode := response.StatusInternalServerError
	html := `<html>
	<head>
		<title>500 Internal Server Error</title>
	</head>
	<body>
		<h1>Internal Server Error</h1>
		<p>Okay, you know what? This one is on me.</p>
	</body>
	</html>`

	writeHTMLResponse(w, statusCode, html)
}

func handlerYourProblem(w *response.Writer, req *request.Request) {
	statusCode := response.StatusBadRequest
	html := `<html>
		<head>
			<title>400 Bad Request</title>
		</head>
		<body>
			<h1>Bad Request</h1>
			<p>Your request honestly kinda sucked.</p>
		</body>
		</html>`

	writeHTMLResponse(w, statusCode, html)
}

func writeHTMLResponse(w *response.Writer, statusCode response.StatusCode, html string) {
	err := w.WriteStatusLine(statusCode)
	if err != nil {
		log.Fatalf("Error writing response status-line: %v", err)
		return
	}

	headers := response.GetDefaultHeaders(len(html))
	headers.SetWithOverride("Content-Type", "text/html")

	err = w.WriteHeaders(headers, false)
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

func handlerProxy(w *response.Writer, req *request.Request) {
	proxyToTarget := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin")
	proxyToUrl := "https://httpbin.org" + proxyToTarget

	fmt.Println(proxyToUrl)

	resp, err := http.Get(proxyToUrl)
	if err != nil {
		// http.Error(w, "Failed to reach target", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	respHeaders := headers.Headers{
		"connection":        "close",
		"content-type":      "text/plain",
		"transfer-encoding": "chunked",
		"trailer":           "X-Content-SHA256, X-Content-Length",
	}

	w.WriteStatusLine(response.StatusOK)
	w.WriteHeaders(respHeaders, false)

	body := []byte{}
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
		body = append(body, buf[:n]...)

	}
	// fmt.Println("End of Body")
	w.WriteChunkedBodyDone(true)

	// calculate body hash
	bodyHash := sha256.Sum256(body)

	// Convert hash to a hex string
	bodyHashHex := hex.EncodeToString(bodyHash[:])

	// Convert length to a string
	bodyLengthStr := strconv.Itoa(len(body))

	trailers := headers.Headers{
		"x-content-sha256": bodyHashHex,
		"x-content-length": bodyLengthStr,
	}

	w.WriteTrailers(trailers)
}

// Reading the whole file into memory (simpler version)
func handlerGetVideo(w *response.Writer, req *request.Request) {
	videoFileBytes, err := os.ReadFile("./assets/vim.mp4")
	if err != nil {
		log.Fatalf("Error reading video file: %v", err)
		return
	}

	err = w.WriteStatusLine(response.StatusOK)
	if err != nil {
		log.Fatalf("Error writing response status-line: %v", err)
		return
	}

	headers := response.GetDefaultHeaders(len(videoFileBytes))
	headers.Set("Content-Type", "video/mp4")

	err = w.WriteHeaders(headers, false)
	if err != nil {
		log.Fatalf("Error writing response headers: %v", err)
		return
	}

	_, err = w.WriteBody(videoFileBytes)
	if err != nil {
		log.Fatalf("Error writing response body: %v", err)
		return
	}
}

// Reading and send the file in chucks with io.Copy (more realistic, browser friendly version)
func handlerGetVideoStream(w *response.Writer, req *request.Request) {
	videoFile, err := os.Open("./assets/vim.mp4")
	if err != nil {
		log.Printf("Error opening video file: %v", err)
		w.WriteStatusLine(response.StatusInternalServerError)
		return
	}
	defer videoFile.Close()

	err = w.WriteStatusLine(response.StatusOK)
	if err != nil {
		log.Printf("Error writing response status-line: %v", err)
		return
	}

	headers := headers.NewHeaders()
	headers.Set("Content-Type", "video/mp4")

	fileInfo, err := videoFile.Stat()
	if err != nil {
		log.Printf("Error getting file info: %v", err)
		w.WriteStatusLine(response.StatusInternalServerError)
		return
	}

	headers.Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))

	err = w.WriteHeaders(headers, false)
	if err != nil {
		log.Printf("Error writing response headers: %v", err)
		return
	}

	// Stream the file instead of loading it all into memory
	_, err = io.Copy(w, videoFile)
	if err != nil {
		// Detect "broken pipe" and avoid logging it as an error
		if strings.Contains(err.Error(), "broken pipe") || strings.Contains(err.Error(), "reset by peer") {
			log.Printf("Client disconnected before full response was sent")
			return
		}
		log.Printf("Error writing response body: %v", err)
	}
}
