package web

import (
	validator "github.com/go-playground/validator/v10"
)

// Validator - go-playground/validator wrapper.
type Validator struct {
	validator *validator.Validate
}

// NewValidator -.
func NewValidator(v *validator.Validate) *Validator {
	if v == nil {
		v = validator.New()
	}

	return &Validator{v}
}

// Validate - go-playground/validator impl.
func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}
