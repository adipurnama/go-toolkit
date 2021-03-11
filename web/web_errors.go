package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/go-playground/validator/v10"
)

// HTTPError -.
type HTTPError struct {
	Code     int         `json:"code"`
	Message  string      `json:"message"`
	Response interface{} `json:"response"`
}

// ErrorDetails ...
type ErrorDetails struct {
	Exception string       `json:"exception"`
	Errors    []ErrorField `json:"errors"`
}

// ErrorField -.
type ErrorField struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *HTTPError) Error() string {
	b, _ := json.Marshal(e)
	return fmt.Sprintf("Error with HTTP StatusCode %d , payload : %s", e.Code, string(b))
}

// NewHTTPValidationError - return new HTTPError caused by validation error.
func NewHTTPValidationError(err error) *HTTPError {
	var fields []ErrorField

	message := errors.Cause(err).Error()

	var validationErrs validator.ValidationErrors

	if ok := errors.As(err, &validationErrs); ok {
		if len(validationErrs) > 0 {
			message = "field validation error found"

			for _, e := range validationErrs {
				fields = append(fields, ErrorField{
					Message: fmt.Sprintf("%s", e),
					Field:   e.Field(),
				})
			}

			if len(fields) > 1 {
				message = "field validation errors found"
			}
		}
	}

	return &HTTPError{
		Code:    http.StatusBadRequest,
		Message: message,
		Response: ErrorDetails{
			Exception: message,
			Errors:    fields,
		},
	}
}
