package main

import "github.com/pkg/errors"

// func maybeError() error {
// 	return nil
// }

func shouldBeError() error {
	return errors.Wrap(guaranteedErr(), "should be error")
}
