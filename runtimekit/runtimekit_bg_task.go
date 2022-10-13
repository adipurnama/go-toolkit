package runtimekit

import (
	"context"
	"runtime/debug"

	"github.com/adipurnama/go-toolkit/log"
	"github.com/pkg/errors"
)

var errRecoverFromPanic = errors.New("recovered from bg task panic")

// ExecuteBackground a function, and log the error when panic found.
// executing a background functions may cause unhandled panic by middlewares.
func ExecuteBackground(fn func()) {
	// Launch a background goroutine.
	go func() {
		// Recover any panic.
		defer func() {
			if r := recover(); r != nil {
				errStack := debug.Stack()

				log.FromCtx(context.Background()).Error(
					errors.Wrapf(errRecoverFromPanic, "caused_by %v", r),
					"found error while executing background task",
					"panic_stack", errStack,
				)
			}
		}()

		// Execute the arbitrary function that we passed as the parameter.
		fn()
	}()
}
