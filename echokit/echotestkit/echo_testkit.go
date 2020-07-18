package echotestkit

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/adipurnama/go-toolkit/web"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// Request return echo.Context and httptest.ResponseRecorder.
func Request(req *http.Request) (echo.Context, *httptest.ResponseRecorder) {
	rec := httptest.NewRecorder()

	e := echo.New()
	e.Validator = web.NewValidator(validator.New())
	ctx := e.NewContext(req, rec)

	return ctx, rec
}

// RequestGET return echo.Context and httptest.ResponseRecorder for GET Request.
func RequestGET(url string) (echo.Context, *httptest.ResponseRecorder) {
	return Request(httptest.NewRequest(http.MethodGet, url, nil))
}

// RequestGETWithParam return echo.Context and httptest.ResponseRecorder for GET Request with URL Param.
func RequestGETWithParam(url string, urlParams map[string]string) (echo.Context, *httptest.ResponseRecorder) {
	ctx, rec := RequestGET(url)

	if urlParams != nil {
		var keys []string

		var values []string

		for key, value := range urlParams {
			keys = append(keys, key)
			values = append(values, value)
		}

		ctx.SetParamNames(keys...)
		ctx.SetParamValues(values...)
	}

	return ctx, rec
}

// RequestPOST return echo.Context and httptest.ResponseRecorder for POST Request.
func RequestPOST(url string, json string) (echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(http.MethodPost, url, strings.NewReader(json))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	return Request(req)
}

// RequestPUT return echo.Context and httptest.ResponseRecorder for POST Request.
func RequestPUT(url string, json string) (ctx echo.Context, rec *httptest.ResponseRecorder) {
	req := httptest.NewRequest(http.MethodPut, url, strings.NewReader(json))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	return Request(req)
}

// RequestDELETE return echo.Context and httptest.ResponseRecorder for DELETE Request.
func RequestDELETE(url string) (ctx echo.Context, rec *httptest.ResponseRecorder) {
	return Request(httptest.NewRequest(http.MethodDelete, url, nil))
}

// RequestDELETEWithParam return echo.Context and httptest.ResponseRecorder for DELETE Request with URL Param.
func RequestDELETEWithParam(url string, urlParams map[string]string) (echo.Context, *httptest.ResponseRecorder) {
	ctx, rec := RequestDELETE(url)

	if urlParams != nil {
		var keys []string

		var values []string

		for key, value := range urlParams {
			keys = append(keys, key)
			values = append(values, value)
		}

		ctx.SetParamNames(keys...)
		ctx.SetParamValues(values...)
	}

	return ctx, rec
}
