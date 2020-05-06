package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// HTTPError -
type HTTPError struct {
	Code       string       `json:"code"`
	Message    string       `json:"message"`
	Exception  string       `json:"exception"`
	StatusCode int          `json:"status_code"`
	Errors     []errorField `json:"errors"`
}

type errorField struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *HTTPError) Error() string {
	b, _ := json.Marshal(e)
	return fmt.Sprintf("Error with HTTP StatusCode %d , payload : %s", e.StatusCode, string(b))
}

// NewHTTPValidationError - return new HTTPError caused by validation error
func NewHTTPValidationError(err error) *HTTPError {
	var fields []errorField
	message := err.Error()

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		if len(validationErrs) > 0 {
			message = "field validation error(s) found"
			for _, e := range validationErrs {
				fields = append(fields, errorField{
					Message: fmt.Sprintf("%s", e),
					Field:   e.Field(),
				})
			}
		}
	}

	return &HTTPError{
		StatusCode: http.StatusBadRequest,
		Code:       "BAD_REQUEST",
		Message:    "Invalid Request",
		Exception:  message,
		Errors:     fields,
	}
}
