package echokit_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adipurnama/go-toolkit/echokit"
	"github.com/adipurnama/go-toolkit/web"
	"github.com/labstack/echo/v4"
)

var (
	errEmptyTraceID   = errors.New("trace ID should not be empty")
	errEmptyRequestID = errors.New("request ID should not be empty")
)

func TestRequestIDMiddleware(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)

	t.Run("new context contains requestID and traceID", func(t *testing.T) {
		checkerHandler := func(ctx echo.Context) error {
			tID := ctx.Request().Header.Get(web.HTTPKeyTraceID)
			if tID == "" {
				return errEmptyTraceID
			}
			rID := ctx.Request().Header.Get(web.HTTPKeyRequestID)
			if rID == "" {
				return errEmptyRequestID
			}
			return nil
		}
		mid := echokit.RequestIDMiddleware()
		handler := mid(checkerHandler)

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if err := handler(c); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("should preserve existing requestID and traceID", func(t *testing.T) {
		traceID := "bhahahah"
		requestID := "mehehehe"
		checkerHandler := func(ctx echo.Context) error {
			tID := ctx.Request().Header.Get(web.HTTPKeyTraceID)
			if tID != traceID {
				return errEmptyTraceID
			}
			rID := ctx.Request().Header.Get(web.HTTPKeyRequestID)
			if rID != requestID {
				return errEmptyRequestID
			}
			return nil
		}
		mid := echokit.RequestIDMiddleware()
		handler := mid(checkerHandler)

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(web.HTTPKeyRequestID, requestID)
		req.Header.Set(web.HTTPKeyTraceID, traceID)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if err := handler(c); err != nil {
			t.Fatal(err)
		}
	})
}
