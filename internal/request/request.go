package request

import (
	"errors"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const (
	CRLF = "\r\n"
)

func RequestFromReader(reader io.Reader) (*Request, error) {
	// buffer := make([]byte, 8, 8)

	entireReqBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	entireReqString := string(entireReqBytes)

	requestParts := strings.Split(entireReqString, CRLF)

	if len(requestParts) == 0 {
		return nil, errors.New("invalid request format")
	}

	requestLine := requestParts[0]

	requestLineParts := strings.Split(requestLine, " ")

	if len(requestLineParts) != 3 {
		return nil, errors.New("invalid request format")
	}

	// Validate fields
	method := requestLineParts[0]
	target := requestLineParts[1]
	httpVersion := requestLineParts[2]

	if httpVersion != "HTTP/1.1" ||
		!isValidHTTPMethod(method) ||
		!strings.HasPrefix(target, "/") {
		return nil, errors.New("invalid request format")
	}

	request := Request{
		RequestLine: RequestLine{
			HttpVersion:   "1.1",
			RequestTarget: target,
			Method:        method,
		},
	}
	return &request, nil
}
