package httpx

import (
	"net/http"

	"github.com/jucardi/go-titan/errors"
)

// Custom HTTP status codes
const (
	// StatusMultipleClientErrors is a custom error code defined in Titan to represent the very rare scenario
	// where multiple client errors occurred while processing the request.
	StatusMultipleClientErrors int = 468

	// StatusMultipleServerErrors is a custom error code defined in Titan to represent the very rare scenario
	// where multiple server errors occurred while processing the request.
	StatusMultipleServerErrors int = 568
)

var (
	statusTextMap = map[int]string{
		StatusMultipleClientErrors: "Multiple Client Errors",
		StatusMultipleServerErrors: "Multiple Server Errors",
	}
)

// StatusText returns a text for the HTTP status code, including custom HTTP status codes registered
// in this package. It returns an empty string if the code is unknown.
func StatusText(code int) string {
	if ret, ok := statusTextMap[code]; ok {
		return ret
	}
	return http.StatusText(code)
}

// RegisterCustomError registers a custom non-existing HTTP status code with it's StatusText. Returns
// an error if the status code is already registered.
func RegisterCustomError(code int, text string) error {
	if text == "" {
		return errors.New("the provided status text cannot be empty")
	}
	if code < 100 || code >= 600 {
		return errors.Format("invalid status code provided (%d), an HTTP status code must be within the range 100-599", code)
	}
	if status := StatusText(code); status != "" {
		return errors.Format("the provided status code (%d) is already registered with the following status text: %s", code, status)
	}
	statusTextMap[code] = text
	return nil
}
