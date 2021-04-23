package v1

import (
	"fmt"

	"github.com/pkg/errors"
)

// EmptyStringValidationError caused by expected `string` field must contains non-empty value.
type EmptyStringValidationError struct {
	FieldName string
}

// Error std error interface impl.
func (e EmptyStringValidationError) Error() string {
	return fmt.Sprintf("validation error: field %s must not be empty", e.FieldName)
}

// Validate grpc_validator.middleware impl.
// will return codes.InvalidArgument.
func (r *HelloRequest) Validate() error {
	if r.Name == "" {
		return EmptyStringValidationError{"HelloRequest.Name"}
	}

	if r.Age <= 0 {
		msg := "HelloRequest.Age must be greater than 0"
		// even though return codes.PermissionDenied, middleware will response with codes.InvalidArgument
		// return status.Errorf(codes.PermissionDenied, msg)

		return errors.Wrap(errors.New(msg), "validation error")
	}

	return nil
}
