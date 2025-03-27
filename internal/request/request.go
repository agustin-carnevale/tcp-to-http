package request

import (
	"errors"
	"fmt"
	"io"

	"github.com/agustin-carnevale/tcp-to-http/internal/headers"
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
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
	REQUEST_PARSING_HEADERS
	REQUEST_PARSING_BODY
	REQUEST_COMPLETED
)

// This is just user for debugging purposes
// var STATE = [4]string{"REQUEST_INITIALIZED", "REQUEST_PARSING_HEADERS", "REQUEST_PARSING_BODY", "REQUEST_COMPLETED"}

const CRLF = "\r\n"
const INITIAL_BUFFER_SIZE = 8

// TODO: case there is a Content-Length but no body at all, implement some check
func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer := make([]byte, INITIAL_BUFFER_SIZE)

	readToIndex := 0
	request := Request{
		state:   REQUEST_INITIALIZED,
		Headers: headers.Headers{},
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
				if _, exists := request.Headers.Get("Content-Length"); exists {
					// if there is a Body, we need to check edge case
					// where we finish reading the request but body < Content-Length
					// in this case we decided to return an error
					// (because no match between actual body and what the header says)
					contentLength, err := contentLengthInt(request.Headers)
					if err != nil {
						fmt.Println(err)
						return nil, err
					}
					if len(request.Body) < contentLength {
						fmt.Println(err)
						return nil, errors.New("body is shorter than Content-Length")
					}
				}
				break
			}
			fmt.Println(err)
			return nil, err
		}

		// update/advance "pointer" after reading n bytes
		readToIndex += numBytesRead

		// PARSE FROM THE BUFFER
		numBytesParsed, err := request.parse(buffer[:readToIndex])
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		// Shift remaining unparsed data to the beginning of the buffer
		copy(buffer, buffer[numBytesParsed:])
		readToIndex -= numBytesParsed
	}

	return &request, nil
}

/*
This parse method will iterate over the data we already read/received until this point
And try to parse as much as possible, until the end or until more data is required
*/
func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.state != REQUEST_COMPLETED {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		totalBytesParsed += n
		// If nothing was parse means we need to receive/read more data in order
		// to have a complete field/part of the request to be parsed
		if n == 0 {
			break
		}
	}
	return totalBytesParsed, nil
}

/*
This parseSingle method will be responsible of parsing only one part/component
of the request at a time. Depending on the state of the parsing, that could be
the entire request-line, a single header (key: value pair), or append the bytes
of the body already received.
*/
func (r *Request) parseSingle(data []byte) (int, error) {

	// if request already completed
	if r.state == REQUEST_COMPLETED {
		return 0, errors.New("error: trying to read data in a done state")
	}

	switch r.state {
	case REQUEST_INITIALIZED:
		// Parse REQUEST-LINE
		// if request is just initialized (first step is parsing request-line)
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
		r.state = REQUEST_PARSING_HEADERS

		return numBytesParsed, nil

	case REQUEST_PARSING_HEADERS:
		// Parse HEADERS
		// if request is done with request-line, start parsing headers
		numBytesParsed, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}

		if done {
			if _, exists := r.Headers.Get("Content-Length"); exists {
				r.state = REQUEST_PARSING_BODY
				return numBytesParsed, nil
			} else {
				r.state = REQUEST_COMPLETED
				return numBytesParsed, nil
			}

		} else {
			if numBytesParsed == 0 {
				// did not parse anything, more data need
				return 0, nil
			} else {
				// it did parse, continued with more headers if any
				return numBytesParsed, nil
			}
		}

	case REQUEST_PARSING_BODY:
		// Parse BODY
		// Append the incoming data to body
		r.Body = append(r.Body, data...)

		contentLength, err := contentLengthInt(r.Headers)
		if err != nil {
			return 0, err
		}

		if len(r.Body) > contentLength {
			return 0, errors.New("body is longer than Content-Length")
		}

		if len(r.Body) == contentLength {
			r.state = REQUEST_COMPLETED
		}

		return len(data), nil

	default:
		return 0, errors.New("error: unknown parser state")
	}
}

func (r *Request) Print() {
	fmt.Println("Request line:")
	fmt.Println("- Method:", r.RequestLine.Method)
	fmt.Println("- Target:", r.RequestLine.RequestTarget)
	fmt.Println("- Version:", r.RequestLine.HttpVersion)
	fmt.Println("")
	fmt.Println("Headers:")
	for key, value := range r.Headers {
		fmt.Printf("- %s: %s\n", key, value)
	}
	if len(r.Body) > 0 {
		fmt.Println("")
		fmt.Println("Body:")
		fmt.Println(string(r.Body))
	}
	// fmt.Println("State:", r.state)
}
