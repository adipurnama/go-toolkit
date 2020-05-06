package echokit

import (
	"context"

	"github.com/adipurnama/go-toolkit/log"
	"github.com/adipurnama/go-toolkit/mask"

	"github.com/adipurnama/go-toolkit/web"

	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

// Middlewares

// RequestIDMiddleware - adds request ID for incoming http request
func RequestIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {

			// trace ID
			tID := ctx.Request().Header.Get(web.HTTPKeyTraceID)
			if tID == "" {
				tID = uuid.NewV4().String()
			}
			ctx.Request().Header.Add(web.HTTPKeyTraceID, tID)
			ctx.Response().Header().Set(web.HTTPKeyTraceID, tID)

			// request ID
			rID := ctx.Request().Header.Get(web.HTTPKeyRequestID)
			if rID == "" {
				rID = uuid.NewV4().String()
			}
			ctx.Request().Header.Add(web.HTTPKeyRequestID, rID)
			ctx.Response().Header().Set(web.HTTPKeyRequestID, rID)

			rCtx := context.WithValue(ctx.Request().Context(), web.ContextKeyRequestID, rID)
			rCtx = context.WithValue(rCtx, web.ContextKeyTraceID, tID)
			ctx.SetRequest(ctx.Request().WithContext(rCtx))
			return next(ctx)
		}
	}
}

// BodyDumpHandler logs incoming request & outgoing response body
func BodyDumpHandler() middleware.BodyDumpHandler {
	return func(ctx echo.Context, reqBody []byte, respBody []byte) {
		// build log payload
		event := log.InfoCtx(ctx.Request().Context()).
			Str("http_method", ctx.Request().Method).
			Str("http_path", mask.URL(ctx.Request().RequestURI)).
			Int("http_status", ctx.Response().Status).
			Str("http_request", mask.URL(string(reqBody)))

		// flush the log message
		event.Str("http_response", mask.URL(string(respBody))).Msg("request completed")
	}
}
