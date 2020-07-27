package helper

import "github.com/pkg/errors"

// GuaranteedErr ...
func GuaranteedErr() error {
	return errors.Wrap(errors.New("guarantee error"), "this function from three.go will absolutely contains error")
}
