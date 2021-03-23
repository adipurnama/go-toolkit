package echokit

import (
	"net/http"

	"github.com/adipurnama/go-toolkit/log"
	"github.com/adipurnama/go-toolkit/web"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func loggerHTTPErrorHandler(w echo.HTTPErrorHandler) echo.HTTPErrorHandler {
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

		var errEchoHTTP *echo.HTTPError
		if ok := errors.As(err, &errEchoHTTP); ok && errEchoHTTP.Internal != nil {
			err = errEchoHTTP.Internal
		}

		if !previouslyCommitted && resp.Committed {
			log.FromCtx(ctx.Request().Context()).Error(err, "request completed with error",
				"path", ctx.Path(),
				"status_code", resp.Status,
			)

			return
		}

		if ctx.Response().Committed {
			return
		}

		code := http.StatusInternalServerError
		msg := "request completed but not handled yet"

		var httpErr *web.HTTPError

		if errEchoHTTP != nil {
			code = errEchoHTTP.Code
			_ = ctx.JSON(errEchoHTTP.Code, errEchoHTTP)
		} else if ok := errors.As(err, &httpErr); ok {
			code = httpErr.Code
			_ = ctx.JSON(httpErr.Code, httpErr)
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
