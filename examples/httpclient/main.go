package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/HereMobilityDevelopers/mediary"
	"github.com/adipurnama/go-toolkit/echokit"
	"github.com/adipurnama/go-toolkit/log"
	"github.com/adipurnama/go-toolkit/web/httpclient"
	echo "github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

const (
	defaultPort = 8081
)

func main() {
	_ = log.NewDevLogger(nil, nil).Set()

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

	ctxClient := httpclient.NewContextHTTPClient(c)
	hClient := dummyAPIClient{client: ctxClient}

	respJSON, httpResp, err := hClient.GetDummy(ctx)
	if err != nil {
		if httpResp != nil {
			log.FromCtx(ctx).Error(err, "got error with response", "status", httpResp.StatusCode)
		} else {
			log.FromCtx(ctx).Error(err, "got error without response")
		}
	} else {
		log.FromCtx(ctx).Info("got response", "resp", respJSON)
	}

	// try to parse resp.Body for the second time
	// should be failed because resp.Body is already closed
	// but we can still get the http statusCode and other info
	var errResp errResponse

	err = json.NewDecoder(httpResp.Body).Decode(&errResp)
	if err != nil {
		err = errors.Wrap(err, "error reading response")
		log.FromCtx(ctx).Error(err, "error parsing second json", "status", httpResp.StatusCode)
	} else {
		log.FromCtx(ctx).Info("got response on second parsing", "resp", respJSON)
	}
}

// ================ Server Handler =============================

func dummy(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, errResponse{
		Status: "OK",
		Code:   http.StatusOK,
	})
}

type errResponse struct {
	Status    string `json:"status"`
	Code      int    `json:"code"`
	Exception string `json:"exception"`
	Message   string `json:"message"`
}

func badRequest(_ echo.Context) error {
	return echo.NewHTTPError(http.StatusBadRequest, errResponse{
		Status:    "Bad Request",
		Code:      http.StatusBadRequest,
		Exception: "bad request is happening",
		Message:   "something bad",
	})
}

func badGateway(_ echo.Context) error {
	return echo.NewHTTPError(http.StatusBadGateway, errResponse{
		Status:    "Bad Gateway",
		Code:      http.StatusBadGateway,
		Exception: "bad gateway is happening",
		Message:   "gateway service error",
	})
}

// =========== API Client =======================

type dummyAPIClient struct {
	client *httpclient.ContextHTTPClient
}

func (c *dummyAPIClient) GetBadRequest(ctx context.Context) (*errResponse, *http.Response, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"http://localhost:8081/bad-request",
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	var errResp errResponse

	resp, err := c.client.Do(ctx, req, &errResp)
	if err != nil {
		return nil, resp, errors.Wrap(err, "failed getting response")
	}

	return &errResp, resp, nil
}

func (c *dummyAPIClient) GetBadGateway(ctx context.Context) (*errResponse, *http.Response, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"http://localhost:8081/bad-gateway",
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	var errResp errResponse

	resp, err := c.client.Do(ctx, req, &errResp)
	if err != nil {
		return nil, resp, errors.Wrap(err, "failed getting response")
	}

	return &errResp, resp, nil
}

func (c *dummyAPIClient) GetDummy(ctx context.Context) (*errResponse, *http.Response, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"http://localhost:8081/dummy",
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	var errResp errResponse

	resp, err := c.client.Do(ctx, req, &errResp)
	if err != nil {
		return nil, resp, errors.Wrap(err, "failed getting response")
	}

	return &errResp, resp, nil
}
