package request

import (
	"errors"
	"strings"
)

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
