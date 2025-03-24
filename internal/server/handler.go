package server

import (
	"io"

	"github.com/agustin-carnevale/tcp-to-http/internal/request"
	"github.com/agustin-carnevale/tcp-to-http/internal/response"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func (h *HandlerError) WriteErrorResponse(w io.Writer) error {

	err := response.WriteStatusLine(w, h.StatusCode)
	if err != nil {
		return err
	}

	contentLength := len(h.Message)
	headers := response.GetDefaultHeaders(contentLength)
	err = response.WriteHeaders(w, headers)
	if err != nil {
		return err
	}

	// response body
	_, err = w.Write([]byte(h.Message))
	if err != nil {
		return err
	}

	return nil
}
