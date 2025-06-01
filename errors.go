package librenms

import (
	"fmt"
	"net/http"
)

type (
	// ErrorResponse represents an error response from the LibreNMS API.
	ErrorResponse struct {
		Response *http.Response `json:"-"`
		Message  string         `json:"message"`
		Status   string         `json:"status"`
	}
)

// Error implements the error interface for ErrorResponse.
func (e *ErrorResponse) Error() string {
	errMsg := fmt.Sprintf("%s %s", e.Response.Status, e.Response.Request.URL.String())
	if e.Message != "" {
		errMsg += fmt.Sprintf(": %s", e.Message)
	}
	return errMsg
}
