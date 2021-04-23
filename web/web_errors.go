package web

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	validator "github.com/go-playground/validator/v10"
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
	return fmt.Sprintf("web.HTTPError code=%d, message=%s", e.Code, e.Message)
}

// NewHTTPValidationError - return new HTTPError caused by validation error.
func NewHTTPValidationError(ctx context.Context, err error) (result *HTTPError) {
	message := errors.Cause(err).Error()
	result = &HTTPError{
		Code:    http.StatusBadRequest,
		Message: message,
	}

	var validationErrs validator.ValidationErrors

	ok := errors.As(err, &validationErrs)
	if !ok {
		return result
	}

	if len(validationErrs) == 0 {
		return result
	}

	var fields []ErrorField

	message = "field validation error found"
	trans := translatorFromContext(ctx)

	for _, e := range validationErrs {
		fieldName := strings.ToLower(e.Field())
		result.Message = e.Translate(trans)

		if trans != nil {
			fields = append(fields, ErrorField{
				Message: e.Translate(trans),
				Field:   fieldName,
			})
		} else {
			fields = append(fields, ErrorField{
				Message: fmt.Sprintf("%s", e),
				Field:   fieldName,
			})
		}
	}

	if len(fields) > 1 {
		message = "field validation errors found"
	}

	result.Response = ErrorDetails{
		Exception: message,
		Errors:    fields,
	}

	return result
}
