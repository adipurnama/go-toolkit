package main

import "github.com/pkg/errors"

func guaranteedErr() error {
	return errors.Wrap(errors.New("guarantee error"), "this function from three.go will absolutely contains error")
}
