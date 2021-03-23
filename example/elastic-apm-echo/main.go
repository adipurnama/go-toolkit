package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/adipurnama/go-toolkit/echokit"
	"github.com/adipurnama/go-toolkit/log"
	"github.com/adipurnama/go-toolkit/tracer"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
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
	errInsufficientMana   ServiceDomainError = errors.New("not enough mana")
	errInvalidRequest     ServiceDomainError = errors.New("invalid request")
	errServiceUnavailable ServiceDomainError = errors.New("service unavailable")
)

type errDataIDNotFound int

func (e errDataIDNotFound) Error() string {
	return fmt.Sprintf("data with ID %d not found", e)
}

func main() {
	if isProdMode {
		_ = log.NewLogger(log.LevelDebug, appName, nil, nil).Set()
	} else {
		_ = log.NewDevLogger(nil, nil).Set()
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

	e.HTTPErrorHandler = errorResponseWriter

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

	var (
		e         ServiceDomainError
		eNotFound *errDataIDNotFound
	)

	if ok := errors.As(err, &e); ok {
		res.Message = e.Error()

		if errors.Is(err, errServiceUnavailable) {
			res.Code = http.StatusBadGateway
		}
	}

	if ok := errors.As(err, &eNotFound); ok {
		res.Code = http.StatusNotFound
		res.Message = eNotFound.Error()
	}

	_ = c.JSON(res.Code, res)
}

type response struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func testErrorHandler(ctx echo.Context) error {
	span := echokit.HandlerSpan(ctx)
	defer span.End()

	err := svcErrFunc(ctx.Request().Context())
	if err != nil {
		return err
	}

	return ctx.String(http.StatusOK, "ok")
}

func testPayment(ctx echo.Context) error {
	span := echokit.HandlerSpan(ctx)
	defer span.End()

	return errors.WithStack(errInsufficientMana)
}

func testSuccess(ctx echo.Context) error {
	span := echokit.HandlerSpan(ctx)
	defer span.End()

	return ctx.String(http.StatusOK, "ok")
}

func testLoginErrorDetailedInfo(ctx echo.Context) error {
	span := echokit.HandlerSpan(ctx)
	defer span.End()

	phone := ctx.QueryParam("phone")
	if strings.ReplaceAll(phone, " ", "") == "" {
		return errors.WithStack(errInvalidRequest)
	}

	if err := svcLogin(ctx.Request().Context(), phone); err != nil {
		msg := fmt.Sprintf("login dengan %s gagal", phone)

		return echo.NewHTTPError(http.StatusBadRequest, msg).
			SetInternal(err)
	}

	return ctx.String(http.StatusAccepted, "Success")
}

func svcLogin(ctx context.Context, phone string) error {
	span := tracer.ServiceFuncSpan(ctx)
	defer span.End()

	// err := errors.WithStack(errServiceUnavailable)
	// err := errServiceUnavailable
	// err := errors.Wrap(errServiceUnavailable, "gateway timeout")

	// return errors.Wrapf(err, "login with phone %s", phone)
	// return errServiceUnavailable
	return errors.Wrapf(errServiceUnavailable, "login with phone %s", phone)
}

func svcErrFunc(ctx context.Context) error {
	span := tracer.ServiceFuncSpan(ctx)
	defer span.End()

	err := repoErrFunc(ctx)

	return err
}

func repoErrFunc(ctx context.Context) error {
	span := tracer.RepositoryFuncSpan(ctx)
	defer span.End()

	uID := 10

	return errors.WithStack(errDataIDNotFound(uID))
}
