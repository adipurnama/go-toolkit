package echokit

import (
	"github.com/adipurnama/go-toolkit/runtimekit"
	"github.com/labstack/echo/v4"
	"go.elastic.co/apm"
)

var skipFuncCount = 2

// HandlerSpan retrieve span for http.Handler / controller type.
func HandlerSpan(ctx echo.Context) *apm.Span {
	tx := apm.TransactionFromContext(ctx.Request().Context())
	return tx.StartSpan(runtimekit.SkippedFunctionName(skipFuncCount), "handler", nil)
}
