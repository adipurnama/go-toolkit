package handler

import (
	"net/http"
	"time"

	"github.com/adipurnama/go-toolkit/log"
	"github.com/labstack/echo/v4"
)

const longProcessTime = 10 * time.Second

// LongOperation simulates long query / api call.
func LongOperation(ctx echo.Context) error {
	rCtx := ctx.Request().Context()

	now := time.Now()

	defer func() {
		log.FromCtx(rCtx).Debug("request ends", "elapsed_time_ms", time.Since(now).Milliseconds())
	}()

	// usually, we call external system with context as parameter
	// e.g. db.QueryRow(ctx, ...), redis.WithContext(ctx).Get(key).Result(), etc..
	// this example using time.Sleep to simulate long running process
	// since it doesn't support ctx param, we use select-case
	select {
	case <-rCtx.Done():
		return echo.NewHTTPError(http.StatusRequestTimeout, "request timeout")
	case <-time.After(longProcessTime):
		return ctx.String(http.StatusOK, "OK. Done.")
	}
}
