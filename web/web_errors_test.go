package web_test

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/pkg/errors"

	"github.com/adipurnama/go-toolkit/web"
	validator "github.com/go-playground/validator/v10"
)

var errValidation = errors.New("some validation error")

type reqParam struct {
	Param1 string `validate:"required"`
	Param2 string `validate:"required"`
}

func TestNewHTTPValidationError(t *testing.T) {
	var r reqParam
	validationErr := validator.New().Struct(r)

	tests := []struct {
		name string
		args error
		want *web.HTTPError
	}{
		{
			"validation error fields",
			validationErr,
			&web.HTTPError{
				Message: "field validation errors found",
				Code:    http.StatusBadRequest,
				Response: web.ErrorDetails{
					Exception: "field validation errors found",
					Errors: []web.ErrorField{
						{
							Field:   "Param1",
							Message: "Key: 'reqParam.Param1' Error:Field validation for 'Param1' failed on the 'required' tag",
						},
						{
							Field:   "Param2",
							Message: "Key: 'reqParam.Param2' Error:Field validation for 'Param2' failed on the 'required' tag",
						},
					},
				},
			},
		},
		{
			"casual error",
			errValidation,
			&web.HTTPError{
				Code:    http.StatusBadRequest,
				Message: errValidation.Error(),
				Response: web.ErrorDetails{
					Exception: "some validation error",
				},
			},
		},
		{
			"wrapped errors",
			errors.Wrap(errValidation, "some additional message"),
			&web.HTTPError{
				Code:    http.StatusBadRequest,
				Message: errValidation.Error(),
				Response: web.ErrorDetails{
					Exception: "some validation error",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := web.NewHTTPValidationError(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHTTPValidationError() = %v, want %v", got, tt.want)
			}
		})
	}
}
