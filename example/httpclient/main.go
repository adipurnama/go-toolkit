package main

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/HereMobilityDevelopers/mediary"
	"github.com/adipurnama/go-toolkit/echokit"
	"github.com/adipurnama/go-toolkit/log"
	"github.com/adipurnama/go-toolkit/web/httpclient"
	"github.com/labstack/echo/v4"
)

const (
	defaultPort = 8081
)

func main() {
	_ = log.NewDevLogger(log.LevelDebug, "sample-httpclient", nil, nil).Set()

	ctx := context.Background()

	e := echo.New()
	e.GET("/dummy", dummy)
	e.GET("/bad-gateway", badGateway)
	e.GET("/bad-request", badRequest)

	go echokit.RunServer(e, &echokit.RuntimeConfig{
		Port: defaultPort,
		Name: "sample-service",
	})

	c := httpclient.NewStdHTTPClient()

	c = mediary.
		Init().
		AddInterceptors(httpclient.LoggerMiddleware).
		WithPreconfiguredClient(c).
		Build()

	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"http://localhost:8081/dummy",
		nil,
	)

	resp, err := c.Do(req)
	if err != nil {
		log.FromCtx(ctx).Error(err, "error getting response")
		return
	}

	resp2, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.FromCtx(ctx).Error(err, "error reading response")

		return
	}

	log.FromCtx(ctx).Info("got response", "json_response", resp2)

	resp.Body.Close()

	req = req.Clone(ctx)
	req.URL.Path = "/bad-request"

	resp, err = c.Do(req)
	if err != nil {
		log.FromCtx(ctx).Error(err, "second request")

		return
	}

	resp.Body.Close()

	req = req.Clone(ctx)
	req.URL.Path = "/bad-gateway"

	resp, err = c.Do(req)
	if err != nil {
		log.FromCtx(ctx).Error(err, "third request")

		return
	}

	resp.Body.Close()
}

func dummy(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, `{"status": "OK"}`)
}

func badRequest(_ echo.Context) error {
	return echo.NewHTTPError(http.StatusBadRequest, "bad request")
}

func badGateway(_ echo.Context) error {
	return echo.NewHTTPError(http.StatusBadGateway, "bad gateway")
}
