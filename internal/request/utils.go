package request

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/agustin-carnevale/tcp-to-http/internal/headers"
)

var validMethods = map[string]struct{}{
	http.MethodGet:     {},
	http.MethodPost:    {},
	http.MethodPut:     {},
	http.MethodPatch:   {},
	http.MethodDelete:  {},
	http.MethodHead:    {},
	http.MethodOptions: {},
	http.MethodTrace:   {},
}

func isValidHTTPMethod(method string) bool {
	_, exists := validMethods[method]
	return exists
}

func contentLengthInt(headers headers.Headers) (int, error) {
	contentLengthString, exists := headers.Get("Content-Length")
	if !exists {
		return 0, errors.New("Content-Length is not defined")
	}

	contentLength, err := strconv.Atoi(contentLengthString)
	if err != nil {
		return 0, errors.New("invalid Content-Length")
	}

	return contentLength, nil
}
