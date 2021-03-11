package echokit

import (
	"net/http"

	"github.com/adipurnama/go-toolkit/errors"
	"github.com/adipurnama/go-toolkit/log"
	"github.com/labstack/echo/v4"
)

// HTTPErrorResponseWriter to write JSON / HTML response for given error type(s).
type HTTPErrorResponseWriter func(err error, ctx echo.Context)

// LoggerHTTPErrorHandler logs error from handler's downstream call.
func LoggerHTTPErrorHandler(w HTTPErrorResponseWriter) echo.HTTPErrorHandler {
	return func(err error, ctx echo.Context) {
		if err == nil {
			return
		}

		previouslyCommitted := ctx.Response().Committed

		if !previouslyCommitted {
			// writer may commit response
			w(errors.Cause(err), ctx)
		}

		resp := ctx.Response()

		if !previouslyCommitted && ctx.Response().Committed {
			log.FromCtx(ctx.Request().Context()).Error(err, "request completed with error",
				"path", ctx.Path(),
				"status_code", resp.Status,
				"content_type", resp.Header().Get("content-type"),
			)

			return
		}

		if !ctx.Response().Committed {
			code := http.StatusInternalServerError
			msg := "request completed but not handled yet"

			var e *echo.HTTPError

			if ok := errors.As(err, &e); ok {
				code = e.Code
				_ = ctx.JSON(e.Code, e)
			} else {
				_ = ctx.JSON(http.StatusInternalServerError, echo.HTTPError{
					Code:     http.StatusInternalServerError,
					Message:  errors.Cause(err).Error(),
					Internal: err,
				})
			}

			if code != http.StatusInternalServerError {
				msg = "http error found"
			}

			log.FromCtx(ctx.Request().Context()).
				Error(err, msg,
					"path", ctx.Path(),
					"status_code", ctx.Response().Status,
				)
		}
	}
}
