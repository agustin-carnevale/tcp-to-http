package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/agustin-carnevale/tcp-to-http/internal/headers"
)

type StatusCode int

const (
	StatusOK        StatusCode = 200
	StatusCreated   StatusCode = 201
	StatusAccepted  StatusCode = 202
	StatusNoContent StatusCode = 204

	StatusBadRequest   StatusCode = 400
	StatusUnauthorized StatusCode = 401
	StatusForbidden    StatusCode = 403
	StatusNotFound     StatusCode = 404

	StatusInternalServerError StatusCode = 500
	StatusNotImplemented      StatusCode = 501
	StatusBadGateway          StatusCode = 502
	StatusServiceUnavailable  StatusCode = 503
)

const CRLF = "\r\n"

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	var statusLine string

	switch statusCode {
	case StatusOK:
		statusLine = "HTTP/1.1 200 OK"
	case StatusBadRequest:
		statusLine = "HTTP/1.1 400 Bad Request"
	case StatusInternalServerError:
		statusLine = "HTTP/1.1 500 Internal Server Error"
	default:
		statusLine = fmt.Sprintf("HTTP/1.1 %d ", statusCode)
	}
	// add '\r\n' at the end of line
	statusLine += CRLF

	_, err := w.Write([]byte(statusLine))
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	return headers.Headers{
		"Content-Length": strconv.Itoa(contentLen),
		"Connection":     "close",
		"Content-Type":   "text/plain",
	}
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	headersString := ""
	for key, value := range headers {
		header := key + ": " + value + CRLF
		headersString += header
	}

	// add '\r\n' at the end of all headers
	headersString += CRLF
	_, err := w.Write([]byte(headersString))

	return err
}
