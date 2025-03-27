package headers

import (
	"bytes"
	"errors"
	"strings"
	"unicode"
)

type Headers map[string]string

const CRLF = "\r\n"

var allowedKeySpecialChars = map[rune]struct{}{
	'!': {}, '#': {}, '$': {}, '%': {}, '&': {},
	'\'': {}, '*': {}, '+': {}, '-': {}, '.': {},
	'^': {}, '_': {}, '`': {}, '|': {}, '~': {},
}

func NewHeaders() Headers {
	return map[string]string{}
}

func (h Headers) Get(key string) (string, bool) {
	key = strings.ToLower(key)
	value, exists := h[key]

	return value, exists
}

func (h Headers) Set(key, value string) {
	// make key lowercase
	key = strings.ToLower(key)

	// if header already exists append new value (comma-separated)
	currentValue, exists := h[key]
	if exists {
		h[key] = currentValue + ", " + value
	} else {
		h[key] = value
	}
}

func (h Headers) SetWithOverride(key, value string) {
	// make key lowercase
	key = strings.ToLower(key)
	h[key] = value
}

// key: value \r\n
func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	endOfHeaderIdx := bytes.Index(data, []byte(CRLF))

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
	h.Set(key, value)

	bytesConsumed := endOfHeaderIdx + len(CRLF)

	// Contemplate edge case is last header (and there is no body)
	CRLFIdx := bytes.Index(data[bytesConsumed:], []byte(CRLF))
	if CRLFIdx == 0 {
		return bytesConsumed, true, nil
	}

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

	if !validHeaderKeyChars(key) {
		return "", "", errors.New("invalid header key")
	}

	//Header Value
	value = strings.TrimSpace(value)

	return key, value, nil
}

func validHeaderKeyChars(key string) bool {
	if len(key) == 0 {
		return false
	}
	for _, r := range key {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			if _, exists := allowedKeySpecialChars[r]; !exists {
				return false
			}
		}
	}
	return true
}
