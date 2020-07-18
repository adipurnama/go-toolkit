package web

import (
	"fmt"
	"net/http"
)

// HTTPClient abstarcts away general http API Call.
type HTTPClient interface {
	Do(r *http.Request) (*http.Response, error)
}

// NewStdHTTPClient returns golang's default httpClient.
func NewStdHTTPClient() *http.Client {
	return http.DefaultClient
}

// HTTPStatusError -.
type HTTPStatusError struct {
	Code int
	Body []byte
}

func (e HTTPStatusError) Error() string {
	return fmt.Sprintf("HTTPStatusCode error status_code=%d response=%s", e.Code, string(e.Body))
}
