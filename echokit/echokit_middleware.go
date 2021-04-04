// Package echokit provides helper for echo web framework app
package echokit

import (
	"context"
	"strings"
	"time"

	"github.com/adipurnama/go-toolkit/log"
	"github.com/lithammer/shortuuid/v3"

	"github.com/adipurnama/go-toolkit/web"

	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"
)

// Middlewares

// RequestIDLoggerMiddleware - adds request ID for incoming http request.
// it also set request with new context with logger
// it's useful when you want to using log package with requestID.
func RequestIDLoggerMiddleware(cfg *RuntimeConfig, o ...Option) echo.MiddlewareFunc {
	opts := options{
		rIDKey: web.HTTPKeyRequestID,
		tIDKey: web.HTTPKeyTraceID,
	}

	for _, o := range o {
		o(&opts)
	}

	m := &rIDLoggerMiddleware{
		rIDKey: opts.rIDKey,
		tIDKey: opts.tIDKey,
		cfg:    cfg,
	}

	return m.handle
}

type rIDLoggerMiddleware struct {
	rIDKey string
	tIDKey string
	cfg    *RuntimeConfig
}

func (m *rIDLoggerMiddleware) handle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pathName := ctx.Request().URL.Path

		if pathName == m.cfg.HealthCheckPath ||
			pathName == m.cfg.InfoCheckPath ||
			pathName == "/metrics" {
			return next(ctx)
		}

		rCtx := log.NewLoggingContext(ctx.Request().Context())
		logger := log.FromCtx(ctx.Request().Context())

		// trace ID
		tID := ctx.Request().Header.Get(m.tIDKey)

		if tID == "" {
			tID = shortuuid.New()
			ctx.Request().Header.Add(web.HTTPKeyTraceID, tID)

			if web.HTTPKeyTraceID != m.tIDKey {
				ctx.Request().Header.Add(m.tIDKey, tID)
			}
		}

		ctx.Response().Header().Set(web.HTTPKeyTraceID, tID)
		logger.AddField("trace_id", tID)

		// request ID
		rID := ctx.Request().Header.Get(m.rIDKey)

		if rID == "" {
			rID = shortuuid.New()

			ctx.Request().Header.Add(web.HTTPKeyRequestID, rID)

			if web.HTTPKeyRequestID != m.rIDKey {
				ctx.Request().Header.Add(m.rIDKey, rID)
			}
		}

		ctx.Response().Header().Set(web.HTTPKeyRequestID, rID)
		logger.AddField("request_id", rID)

		rCtx = log.AddToContext(rCtx, logger)

		ctx.SetRequest(ctx.Request().WithContext(rCtx))

		return next(ctx)
	}
}

type options struct {
	rIDKey string
	tIDKey string
}

// Option sets options for request middleware.
type Option func(*options)

// WithRequestIDKey returns an Option which sets `key` as request-ID lookup
// to use for logging server requests.
func WithRequestIDKey(key string) Option {
	return func(o *options) {
		if key != "" {
			o.rIDKey = key
		}
	}
}

// WithTraceIDKey returns an Option which sets `key` as trace-ID lookup
// to use for logging server requests.
func WithTraceIDKey(key string) Option {
	return func(o *options) {
		if key != "" {
			o.tIDKey = key
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

// TimeoutConfig request timeout configuration
// default value:
//  * timeout: 7 seconds
//	* middleware.DefaultSkipper / apply to all url
type TimeoutConfig struct {
	Timeout time.Duration
	Skipper middleware.Skipper
}

// TimeoutMiddleware sets upstream request context's timeout.
func TimeoutMiddleware(cfg *TimeoutConfig) echo.MiddlewareFunc {
	// setup default value
	if cfg.Skipper == nil {
		cfg.Skipper = middleware.DefaultSkipper
	}

	if cfg.Timeout <= 0 {
		cfg.Timeout = defaultReqTimeout
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if skip := cfg.Skipper(ctx); skip {
				return next(ctx)
			}

			req := ctx.Request()

			rCtx, cancel := context.WithTimeout(req.Context(), cfg.Timeout)
			defer cancel()

			ctx.SetRequest(req.WithContext(rCtx))

			return next(ctx)
		}
	}
}
