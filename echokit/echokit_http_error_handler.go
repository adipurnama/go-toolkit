package echokit

import (
	"net/http"

	"github.com/adipurnama/go-toolkit/log"
	"github.com/adipurnama/go-toolkit/web"
	echo "github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type defaultErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func loggerHTTPErrorHandler(w echo.HTTPErrorHandler) echo.HTTPErrorHandler {
	return func(err error, ctx echo.Context) {
		logger := log.FromCtx(ctx.Request().Context())
		msg := "request completed with error"

		prevCommitted := ctx.Response().Committed

		if !ctx.Response().Committed {
			// writer may commit response
			w(err, ctx)
		}

		if ctx.Response().Committed {
			if !prevCommitted {
				logErrorAndResponse(logger, msg, err, ctx)
			}

			return
		}

		// found error & response not yet written

		// check for echo.NewHTTPError returned from handler / controller
		var errEchoHTTP *echo.HTTPError
		if ok := errors.As(err, &errEchoHTTP); ok {
			if errEchoHTTP.Internal != nil {
				err = errEchoHTTP.Internal
			}

			errWriteResp := ctx.JSON(errEchoHTTP.Code, errEchoHTTP)

			if errWriteResp != nil {
				logger.Error(errWriteResp, "error writing JSON response", "path", ctx.Request().URL.Path)
			}

			logErrorAndResponse(logger, msg, err, ctx)

			return
		}

		// check for web.Validation error
		var httpErr *web.HTTPError
		if ok := errors.As(err, &httpErr); ok {
			errWriteResp := ctx.JSON(httpErr.Code, httpErr)

			if errWriteResp != nil {
				logger.Error(errWriteResp, "error writing JSON response", "path", ctx.Request().URL.Path)
			}

			logErrorAndResponse(logger, msg, err, ctx)

			return
		}

		// unhandled errors returned types from controller / handler
		resp := defaultErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
		}

		errWriteResp := ctx.JSON(resp.Code, resp)
		if errWriteResp != nil {
			logger.Error(errWriteResp, "error writing JSON response", "path", ctx.Request().URL.Path)
		}

		logErrorAndResponse(logger, "request completed with unhandled error. add error type inspection in your echo.HTTPErrorHandler", err, ctx)
	}
}

func logErrorAndResponse(l *log.Logger, msg string, err error, ctx echo.Context) {
	if ctx.Response().Status >= http.StatusInternalServerError {
		l.Error(err, msg,
			"path", ctx.Request().URL.Path,
			"status_code", ctx.Response().Status,
		)
	} else {
		l.Info(msg,
			"error", err,
			"path", ctx.Request().URL.Path,
			"status_code", ctx.Response().Status,
		)
	}
}
