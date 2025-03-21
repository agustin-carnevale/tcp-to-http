package request

import (
	"errors"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	state       requestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type requestState int

const (
	REQUEST_INITIALIZED requestState = iota
	REQUEST_COMPLETED
)

const CRLF = "\r\n"
const INITIAL_BUFFER_SIZE = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer := make([]byte, INITIAL_BUFFER_SIZE)

	readToIndex := 0
	request := Request{
		state: REQUEST_INITIALIZED,
	}

	for request.state != REQUEST_COMPLETED {
		if readToIndex >= len(buffer) {
			// if buffer is full duplicate size/capacity
			newBuffer := make([]byte, len(buffer)*2)
			copy(newBuffer, buffer)
			buffer = newBuffer
		}

		// READ INTO BUFFER
		numBytesRead, err := reader.Read(buffer[readToIndex:])
		if err != nil {
			if err == io.EOF {
				request.state = REQUEST_COMPLETED
				break
			}
			return nil, err
		}

		// update/advance "pointer" after reading n bytes
		readToIndex += numBytesRead

		// PARSE FROM THE BUFFER
		numBytesParsed, err := request.parse(buffer[:readToIndex])
		if err != nil {
			return nil, err
		}

		// Shift remaining unparsed data to the beginning of the buffer
		copy(buffer, buffer[numBytesParsed:])
		readToIndex -= numBytesParsed

	}
	return &request, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	entireReqString := string(data)
	requestParts := strings.Split(entireReqString, CRLF)

	// Validate parts divided by CRLF
	if len(requestParts) < 2 {
		// means not enough data to read entire request-line
		return nil, 0, nil
	}

	// First part should be a RequestLine
	requestLineString := requestParts[0]

	requestLine, err := parseRequestLineFromString(requestLineString)
	if err != nil {
		return nil, 0, err
	}

	// if parse was successful then I parsed the whole requestLine string + the CRLF (which is ignored)
	bytesParsed := len(requestLineString) + len(CRLF)

	return requestLine, bytesParsed, nil
}

func parseRequestLineFromString(requestLineString string) (*RequestLine, error) {
	requestLineParts := strings.Split(requestLineString, " ")

	// Validate RequestLine parts
	if len(requestLineParts) != 3 {
		return nil, errors.New("invalid request format")
	}

	// RequestLine fields
	method := requestLineParts[0]
	target := requestLineParts[1]
	version := requestLineParts[2]

	//Validate version
	versionParts := strings.Split(version, "/")
	if len(versionParts) != 2 {
		return nil, errors.New("invalid request http version format")
	}
	httpVersion := versionParts[1]
	if httpVersion != "1.1" {
		return nil, errors.New("invalid request version")
	}

	//Validate method
	if !isValidHTTPMethod(method) {
		return nil, errors.New("invalid request method")
	}

	//Validate target
	if !strings.HasPrefix(target, "/") {
		return nil, errors.New("invalid request target")
	}

	requestLine := RequestLine{
		HttpVersion:   httpVersion,
		RequestTarget: target,
		Method:        method,
	}

	return &requestLine, nil
}

func (r *Request) parse(data []byte) (int, error) {

	// if request already completed
	if r.state == REQUEST_COMPLETED {
		return 0, errors.New("error: trying to read data in a done state")
	}

	// if request is not completed
	if r.state == REQUEST_INITIALIZED {
		requestLine, numBytesParsed, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if numBytesParsed == 0 {
			// couldn't parse yet, more data needed
			return 0, nil
		}

		// if bytes consumed then update requestLine and state
		r.RequestLine = *requestLine
		r.state = REQUEST_COMPLETED

		return numBytesParsed, nil
	}

	return 0, errors.New("error: unknown parser state")
}
