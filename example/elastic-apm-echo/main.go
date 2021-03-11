package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/adipurnama/go-toolkit/echokit"
	"github.com/adipurnama/go-toolkit/errors"
	"github.com/adipurnama/go-toolkit/log"
	"github.com/adipurnama/go-toolkit/runtimekit"
	"github.com/labstack/echo/v4"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmechov4"
)

const (
	isProdMode = false
	port       = 8022
	appName    = "elastic-apm-echo-samplekit"
)

// ServiceDomainError for service specific error types.
type ServiceDomainError error

var (
	errDataNotFound       ServiceDomainError = errors.New("data not found")
	errInsufficientMana   ServiceDomainError = errors.New("not enough mana")
	errInvalidRequest     ServiceDomainError = errors.New("invalid request")
	errServiceUnavailable ServiceDomainError = errors.New("service unavailable")
)

func main() {
	if isProdMode {
		_ = log.NewLogger(log.LevelDebug, appName, nil, nil).Set()
	} else {
		_ = log.NewDevLogger(log.LevelDebug, appName, nil, nil).Set()
	}

	eCfg := echokit.RuntimeConfig{
		Port:            port,
		Name:            appName,
		HealthCheckPath: "/health/info",
		HealthCheckFunc: func(ctx context.Context) error {
			return nil
		},
	}

	e := echo.New()
	e.Use(
		echokit.RequestIDLoggerMiddleware(&eCfg),
		apmechov4.Middleware(),
	)

	e.HTTPErrorHandler = echokit.LoggerHTTPErrorHandler(errorResponseWriter)

	e.GET("/errors", testErrorHandler)
	e.GET("/success", testSuccess)
	e.GET("/payment", testPayment)
	e.GET("/login", testLoginErrorDetailedInfo)

	echokit.RunServer(e, &eCfg)
}

func errorResponseWriter(err error, c echo.Context) {
	res := response{
		Code:    http.StatusBadRequest,
		Message: "error nih ye",
	}

	var e ServiceDomainError

	if ok := errors.As(err, &e); ok {
		switch e {
		case errDataNotFound:
			_ = c.JSON(res.Code, res)
		case errInsufficientMana:
			res.Message = e.Error()
			_ = c.JSON(http.StatusBadRequest, res)
		case errServiceUnavailable:
			res.Code = http.StatusBadGateway
			_ = c.JSON(res.Code, res)
		}
	}
}

type response struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func testErrorHandler(ctx echo.Context) error {
	tx := apm.TransactionFromContext(ctx.Request().Context())

	span := tx.StartSpan(runtimekit.CallerName(), "handler", nil)
	defer span.End()

	err := svcErrFunc(ctx.Request().Context())
	if err != nil {
		return errors.WrapFunc(err)
	}

	return ctx.String(http.StatusOK, "ok")
}

func testPayment(ctx echo.Context) error {
	tx := apm.TransactionFromContext(ctx.Request().Context())

	span := tx.StartSpan(runtimekit.CallerName(), "handler", nil)
	defer span.End()

	return errInsufficientMana
}

func testSuccess(ctx echo.Context) error {
	tx := apm.TransactionFromContext(ctx.Request().Context())

	span := tx.StartSpan(runtimekit.CallerName(), "handler", nil)
	defer span.End()

	return ctx.String(http.StatusOK, "ok")
}

func testLoginErrorDetailedInfo(ctx echo.Context) error {
	tx := apm.TransactionFromContext(ctx.Request().Context())

	span := tx.StartSpan(runtimekit.CallerName(), "handler", nil)
	defer span.End()

	phone := ctx.QueryParam("phone")
	if strings.ReplaceAll(phone, " ", "") == "" {
		return errInvalidRequest
	}

	if err := svcLogin(ctx.Request().Context(), phone); err != nil {
		msg := fmt.Sprintf("login dengan %s gagal", phone)

		return echo.NewHTTPError(http.StatusBadRequest, msg).
			SetInternal(err)
	}

	return ctx.String(http.StatusAccepted, "Success")
}

func svcLogin(_ context.Context, phone string) error {
	return errors.WrapFunc(errServiceUnavailable, "phone", phone)
}

func svcErrFunc(ctx context.Context) error {
	tx := apm.TransactionFromContext(ctx)

	span := tx.StartSpan(runtimekit.CallerName(), "service", nil)
	defer span.End()

	err := repoErrFunc(ctx)

	return errors.WrapFuncMsg(err, "got error from repo")
}

func repoErrFunc(ctx context.Context) error {
	tx := apm.TransactionFromContext(ctx)

	span := tx.StartSpan(runtimekit.CallerName(), "repo", nil)
	defer span.End()

	return errors.WrapFunc(errDataNotFound)
}
