package request

import "net/http"

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
