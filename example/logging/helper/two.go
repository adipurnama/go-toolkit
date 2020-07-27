package helper

import "github.com/pkg/errors"

// ShouldBeError ...
func ShouldBeError() error {
	return errors.Wrap(GuaranteedErr(), "should be error")
}
