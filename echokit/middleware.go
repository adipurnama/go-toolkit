package echokit

import (
	"github.com/adipurnama/go-toolkit/log"
	"github.com/adipurnama/go-toolkit/mask"
	uuid "github.com/satori/go.uuid"

	"github.com/adipurnama/go-toolkit/web"

	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"
)

// Middlewares

// RequestIDLoggerMiddleware - adds request ID for incoming http request.
// it also set request with new context with logger
// it's useful when you want to using log package with requestID
func RequestIDLoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			logger := log.FromCtx(ctx.Request().Context())

			// trace ID
			tID := ctx.Request().Header.Get(web.HTTPKeyTraceID)

			if tID == "" {
				tID = uuid.NewV4().String()
				ctx.Request().Header.Add(web.HTTPKeyTraceID, tID)
			}

			ctx.Response().Header().Set(web.HTTPKeyTraceID, tID)
			logger.AddField("trace_id", web.ContextKeyTraceID)

			// request ID
			rID := ctx.Request().Header.Get(web.HTTPKeyRequestID)

			if rID == "" {
				rID = uuid.NewV4().String()
				ctx.Request().Header.Add(web.HTTPKeyRequestID, rID)
			}

			ctx.Response().Header().Set(web.HTTPKeyRequestID, rID)
			logger.AddField("request_id", web.HTTPKeyRequestID)

			rCtx := log.NewContextLogger(ctx.Request().Context(), logger)
			ctx.SetRequest(ctx.Request().WithContext(rCtx))

			return next(ctx)
		}
	}
}

// BodyDumpHandler logs incoming request & outgoing response body.
func BodyDumpHandler() middleware.BodyDumpHandler {
	return func(ctx echo.Context, reqBody []byte, respBody []byte) {
		log.FromCtx(ctx.Request().Context()).Info(
			"request completed",
			"http_method", ctx.Request().Method,
			"http_path", mask.URL(ctx.Request().RequestURI),
			"http_status", ctx.Response().Status,
			"http_request", mask.URL(string(reqBody)),
			"http_response", mask.URL(string(respBody)),
		)
	}
}
