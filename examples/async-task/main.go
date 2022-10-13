package main

import (
	// stdLog "log"
	"net/http"
	"time"

	"github.com/adipurnama/go-toolkit/log"
	"github.com/adipurnama/go-toolkit/runtimekit"
	"github.com/labstack/echo/v4"
)

func main() {
	log.NewDevLogger(nil, nil).Set()

	e := echo.New()
	e.POST("/user-registration", userRegistrationHandler)
	e.Start(":8081")
}

func userRegistrationHandler(ctx echo.Context) error {
	// send registration email in async/background
	runtimekit.ExecuteBackground(func() {
		log.FromCtx(ctx.Request().Context()).Info("sending email...")
		time.Sleep(5 * time.Second)

		// Uncomment to simulate panic / unexpected error
		// stdLog.Panic("SIMULATES ERROR HERE..")

		log.FromCtx(ctx.Request().Context()).Info("email sent")
	})

	return ctx.String(http.StatusAccepted, http.StatusText(http.StatusAccepted))
}
