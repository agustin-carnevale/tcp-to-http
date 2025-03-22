package headers

import (
	"bytes"
	"errors"
	"strings"

	"github.com/agustin-carnevale/tcp-to-http/internal/request"
)

type Headers map[string]string

// key: value \r\n
func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	endOfHeaderIdx := bytes.Index(data, []byte(request.CRLF))

	// More data needed
	if endOfHeaderIdx == -1 {
		return 0, false, nil
	}
	// If data starts with CRLF means we are done with the Headers
	if endOfHeaderIdx == 0 {
		return 2, true, nil
	}

	// parse header
	key, value, err := parseHeader(string(data[:endOfHeaderIdx]))
	if err != nil {
		return 0, false, err
	}

	// store header
	h[key] = value

	bytesConsumed := endOfHeaderIdx + len(request.CRLF)
	return bytesConsumed, false, nil
}

func parseHeader(header string) (key string, value string, err error) {
	key, value, found := strings.Cut(header, ":")

	if !found {
		return "", "", errors.New("invalid header format")
	}

	// Header key
	// Trim whitespace at the beginig
	key = strings.TrimLeft(key, " ")

	// Check there is no space the end
	// between key and : (this "key  : value" is not valid)
	if key != strings.TrimSpace(key) {
		return "", "", errors.New("invalid header format (space between key and :)")
	}

	//Header Value
	value = strings.TrimSpace(value)

	return key, value, nil
}
