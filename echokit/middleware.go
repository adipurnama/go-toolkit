// Package echokit provides helper for echo web framework app
package echokit

import (
	"strings"

	"github.com/adipurnama/go-toolkit/log"
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
			rCtx := log.NewLoggingContext(ctx.Request().Context())
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

			ctx.SetRequest(ctx.Request().WithContext(rCtx))

			return next(ctx)
		}
	}
}

// BodyDumpHandler logs incoming request & outgoing response body.
func BodyDumpHandler(skipper middleware.Skipper) echo.MiddlewareFunc {
	return middleware.BodyDumpWithConfig(
		middleware.BodyDumpConfig{
			Skipper: skipper,
			Handler: bodyDumpHandlerFunc,
		},
	)
}

func bodyDumpHandlerFunc(c echo.Context, reqBody []byte, respBody []byte) {
	l := log.FromCtx(c.Request().Context())
	respStr := strings.Replace(string(respBody), "}\n", "}", 1)
	reqStr := strings.Replace(string(reqBody), "}\n", "}", 1)

	l.Debug("REST request completed",
		"http.response", respStr,
		"http.request", reqStr,
		"http.status_code", c.Response().Status,
		"http.header", c.Request().Header,
	)
}
