package echokit

import (
	"fmt"
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/adipurnama/go-toolkit/log"
	"github.com/adipurnama/go-toolkit/web"
)

type defaultErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func loggerHTTPErrorHandler(w echo.HTTPErrorHandler) echo.HTTPErrorHandler {
	return func(err error, ctx echo.Context) {
		logger := log.FromCtx(ctx.Request().Context())
		msg := fmt.Sprintf("%s %s - request completed with error", ctx.Request().Method, ctx.Request().URL.Path)

		prevCommitted := ctx.Response().Committed

		if !prevCommitted {
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

		// check for gRPC API call returns
		// convert it to *echo.HTTPError
		st, ok := status.FromError(errors.Cause(err))
		if ok {
			code := grpcToHTTPStatusCode(st.Code())
			err = echo.NewHTTPError(code, st.Message()).SetInternal(err)
		}

		// check for echo.NewHTTPError returned from handler / controller
		var errEchoHTTP *echo.HTTPError
		if ok := errors.As(err, &errEchoHTTP); ok {
			if errEchoHTTP.Internal != nil {
				err = errEchoHTTP.Internal
			}

			errWriteResp := ctx.JSON(errEchoHTTP.Code, errEchoHTTP)

			if errWriteResp != nil {
				logger.WarnError(errWriteResp, "error writing JSON response", "path", ctx.Request().URL.Path)
			}

			logErrorAndResponse(logger, msg, err, ctx)

			return
		}

		// check for web.Validation error
		var httpErr *web.HTTPError
		if ok := errors.As(err, &httpErr); ok {
			errWriteResp := ctx.JSON(httpErr.Code, httpErr)

			if errWriteResp != nil {
				logger.WarnError(errWriteResp, "error writing JSON response", "path", ctx.Request().URL.Path)
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
			logger.WarnError(errWriteResp, "error writing JSON response", "path", ctx.Request().URL.Path)
		}

		logErrorAndResponse(
			logger,
			"request completed with unhandled error. add error type inspection in your echo.HTTPErrorHandler",
			err,
			ctx,
		)
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

// non-standard nginx response for client-cancelled operation.
const httpStatusCancelled = 499

// taken from https://github.com/cloudendpoints/esp/blob/d144b204c5fa380e8ccb0dc9ba33ea32fdef8871/src/api_manager/utils/status.cc#L289
func grpcToHTTPStatusCode(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return httpStatusCancelled
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.FailedPrecondition:
		return http.StatusPreconditionFailed
	case codes.Aborted:
		return http.StatusPreconditionRequired
	case codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
