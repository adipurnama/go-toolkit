package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

var errInternal = errors.New("internal handler error")

// Success ...
// GET /success.
func Success() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "success")
	}
}

// EchoHTTPError ...
// GET /echo-error/:code.
func EchoHTTPError() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		code := ctx.Param("code")

		codeInt, err := strconv.Atoi(code)
		if err != nil {
			return err
		}

		return echo.NewHTTPError(codeInt, "expected error")
	}
}

// InternalError expects 500 error returned.
func InternalError() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return errInternal
	}
}

// Panic simulates panic handler.
func Panic() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		log.Panic("I'd like to PANIC!!! Please.")

		return nil
	}
}
