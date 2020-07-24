package helper

import "github.com/pkg/errors"

// func maybeError() error {
// 	return nil
// }

func ShouldBeError() error {
	return errors.Wrap(GuaranteedErr(), "should be error")
}
