package helper

import "github.com/pkg/errors"

// DefinitelyError ...
func DefinitelyError() error {
	return errors.Wrap(ShouldBeError(), "definitelyError")
}
