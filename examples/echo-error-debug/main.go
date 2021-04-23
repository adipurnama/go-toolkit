package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/adipurnama/go-toolkit/echokit"
	"github.com/adipurnama/go-toolkit/examples/echo-error-debug/handler"
	"github.com/adipurnama/go-toolkit/examples/echo-error-debug/middleware"
	echo "github.com/labstack/echo/v4"
	// echo_mdw "github.com/labstack/echo/v4/middleware".
)

func main() {
	e := echo.New()
	e.HideBanner = true
	e.Use(
		// echo_mdw.Recover(),
		middleware.DummyErrorReportingMiddleware,
		middleware.ErrorLoggerMiddleware,
	)

	e.GET("/errors", handler.InternalError())
	e.GET("/echo-error/:code", handler.EchoHTTPError())
	e.GET("/success", handler.Success())
	e.GET("/panic", handler.Panic())

	e.HTTPErrorHandler = errorHandler

	echokit.PrintRoutes(e)

	if err := e.Start(":8080"); err != nil {
		log.Println(err)
	}
}

type response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func errorHandler(err error, ctx echo.Context) {
	if ctx.Response().Committed {
		log.Printf("errorHandler: response_already_committed got error %s. Abort.", err)
		return
	}

	var (
		resp        response
		errEchoHTTP *echo.HTTPError
	)

	resp.Code = http.StatusInternalServerError
	resp.Message = http.StatusText(resp.Code)

	if ok := errors.As(err, &errEchoHTTP); ok {
		resp.Code = errEchoHTTP.Code
		resp.Message = fmt.Sprintf("%s", errEchoHTTP.Message)
	}

	if writeErr := ctx.JSON(resp.Code, resp); writeErr != nil {
		log.Println("errorHandler: error writing response:", writeErr)
	}

	log.Println("errorHandler: done writing response for error:", err)
}
