package middleware

import (
	"fmt"
	"log"
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

var errRecoveredPanic = errors.New("recovered from panic")

// ErrorLoggerMiddleware ...
func ErrorLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		err := next(ctx)
		if err != nil {
			log.Println("inside ErrorLoggerMiddleware, found error:", err)
		}

		return err
	}
}

// DummyErrorReportingMiddleware ...
func DummyErrorReportingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var handlerErr error

		defer func() {
			if v := recover(); v != nil {
				err, ok := v.(error)
				if !ok {
					err = errors.Wrap(errRecoveredPanic, fmt.Sprint(v))
				}

				ctx.Error(err)

				log.Println("reporting panic error:", err)
			}

			if handlerErr == nil {
				log.Println("reporter: no error found")

				return
			}

			if !ctx.Response().Committed {
				ctx.Error(handlerErr)
			}

			if ctx.Response().Status >= http.StatusInternalServerError {
				log.Println("reporting handler error:", handlerErr,
					"status", ctx.Response().Status)
			}
		}()

		resp := ctx.Response()

		handlerErr = next(ctx)
		if handlerErr != nil {
			log.Println("inside dummy reporter, found error: ", handlerErr)

			resp.Status = http.StatusInternalServerError

			var echoErr *echo.HTTPError

			if ok := errors.As(handlerErr, &echoErr); ok {
				log.Println("inside dummy reporter, found echo error: ", echoErr)

				resp.Status = echoErr.Code
			}

			return handlerErr
		}

		if !resp.Committed {
			log.Println("response not committed, writing OK header")
			resp.WriteHeader(http.StatusOK)

			return handlerErr
		}

		log.Println("response already committed:", resp.Committed,
			"status:", resp.Status)

		return handlerErr
	}
}
