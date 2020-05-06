package main

import "github.com/pkg/errors"

func definitelyError() error {
	return errors.Wrap(shouldBeError(), "definitelyError")
}
