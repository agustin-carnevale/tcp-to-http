package response

import (
	"fmt"
	"net"
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

const (
	WriteStatusLine WriterState = iota
	WriteHeaders
	WriteBody
)

type WriterState int

type Writer struct {
	Connection net.Conn
	state      WriterState
}

func (w *Writer) Write(data []byte) (int, error) {
	return w.Connection.Write(data)
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != WriteStatusLine {
		return fmt.Errorf("cannot write status line in state %d", w.state)
	}
	defer func() { w.state = WriteHeaders }()

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
		"content-length": strconv.Itoa(contentLen),
		"connection":     "close",
		"content-type":   "text/plain",
	}
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != WriteHeaders {
		return fmt.Errorf("cannot write headers in state %d", w.state)
	}
	defer func() { w.state = WriteBody }()

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

func (w *Writer) WriteBody(body []byte) (int, error) {
	if w.state != WriteBody {
		return 0, fmt.Errorf("cannot write body in state %d", w.state)
	}
	return w.Write(body)
}
