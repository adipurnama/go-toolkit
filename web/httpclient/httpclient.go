// Package httpclient is http.Client helpers
package httpclient

import (
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/HereMobilityDevelopers/mediary"
	"github.com/adipurnama/go-toolkit/log"
)

const (
	defaultTimeout       = 5 * time.Second
	defaultMaxConnection = 50
)

type options struct {
	maxConnection int
	timeout       time.Duration
}

// Option sets options for http client.
type Option func(*options)

// NewStdHTTPClient returns golang's default httpClient.
func NewStdHTTPClient(opts ...Option) *http.Client {
	o := options{
		maxConnection: defaultMaxConnection,
		timeout:       defaultTimeout,
	}

	for _, opt := range opts {
		opt(&o)
	}

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = o.maxConnection
	t.MaxConnsPerHost = o.maxConnection
	t.MaxIdleConnsPerHost = o.maxConnection

	return &http.Client{
		Timeout:   o.timeout,
		Transport: t,
	}
}

// WithMaxConnection returns an Option which sets http client connection pool size
// default is 50.
func WithMaxConnection(m int) Option {
	return func(o *options) {
		if m > 0 {
			o.maxConnection = m
		}
	}
}

// WithTimeout returns an Option which sets http client's timeout
// default is 5 seconds.
func WithTimeout(t time.Duration) Option {
	return func(o *options) {
		if t.Seconds() > 0 {
			o.timeout = t
		}
	}
}

// HTTP Client mediary middleware

// LoggerMiddleware logs request and response of http calls.
func LoggerMiddleware(req *http.Request, handler mediary.Handler) (*http.Response, error) {
	logger := log.FromCtx(req.Context())

	keyVals := []interface{}{"url", req.URL.String(), "method", req.Method}

	dumpReq, err := httputil.DumpRequest(req, true)

	if err != nil {
		keyVals = append(keyVals, "req_headers", req.Header)
	} else {
		keyVals = append(keyVals, "req_body", dumpReq)
	}

	resp, err := handler(req)
	if err != nil {
		return resp, err
	}

	if rBody, err2 := httputil.DumpResponse(resp, true); err2 == nil {
		keyVals = append(keyVals, "resp_body", rBody)
	} else {
		keyVals = append(keyVals, "resp_body", "err reading resp_body")
	}

	logger.Debug("http.Client call completed", keyVals...)

	return resp, err
}
