package helper

import (
	"github.com/adipurnama/go-toolkit/errors"
)

// DefinitelyError ...
func DefinitelyError() error {
	return errors.WrapFunc(ShouldBeError())
}
