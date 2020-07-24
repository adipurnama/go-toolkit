package helper

import "github.com/pkg/errors"

func DefinitelyError() error {
	return errors.Wrap(ShouldBeError(), "definitelyError")
}
