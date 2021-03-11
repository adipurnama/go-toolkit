package helper

import (
	"github.com/adipurnama/go-toolkit/errors"
)

// ShouldBeError ...
func ShouldBeError() error {
	return errors.WrapFunc(GuaranteedErr(), "info", "two.go")
}
