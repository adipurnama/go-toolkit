package echoapmkit

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"

	"github.com/adipurnama/go-toolkit/runtimekit"
	echo "github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"go.elastic.co/apm"
	apmhttp "go.elastic.co/apm/module/apmhttp"
)

const skipFuncCount = 2

// HandlerSpan retrieve span for http.Handler / controller type.
func HandlerSpan(ctx echo.Context) *apm.Span {
	tx := apm.TransactionFromContext(ctx.Request().Context())
	return tx.StartSpan(runtimekit.SkippedFunctionName(skipFuncCount), "http.handler", nil)
}

// MODIFIED from https://github.com/elastic/apm-agent-go/blob/master/module/apmechov4/middleware.go

// ElasticAPMMiddleware returns a new Echo middleware handler for tracing
// requests and reporting errors.
//
// This middleware will recover and report panics, so it can
// be used instead of echo/middleware.Recover.
//
// By default, the middleware will use apm.DefaultTracer.
// Use WithTracer to specify an alternative tracer.
func ElasticAPMMiddleware(o ...APMOption) echo.MiddlewareFunc {
	opts := options{
		tracer:         apm.DefaultTracer,
		requestIgnorer: apmhttp.DefaultServerRequestIgnorer(),
	}

	for _, o := range o {
		o(&opts)
	}

	return func(h echo.HandlerFunc) echo.HandlerFunc {
		m := &middleware{
			tracer:         opts.tracer,
			handler:        h,
			requestIgnorer: opts.requestIgnorer,
		}

		return m.handle
	}
}

type middleware struct {
	handler        echo.HandlerFunc
	tracer         *apm.Tracer
	requestIgnorer apmhttp.RequestIgnorerFunc
}

var errPanicInternal = errors.New("found panic while serving request")

func (m *middleware) handle(c echo.Context) error {
	req := c.Request()
	if !m.tracer.Recording() || m.requestIgnorer(req) {
		return m.handler(c)
	}

	name := req.Method + " " + c.Path()

	tx, req := apmhttp.StartTransaction(m.tracer, name, req)
	defer tx.End()

	c.SetRequest(req)

	body := m.tracer.CaptureHTTPRequestBody(req)
	resp := c.Response()

	resp.Status = http.StatusInternalServerError

	var handlerErr error

	defer func() {
		if v := recover(); v != nil {
			err, ok := v.(error)
			if !ok {
				err = errors.Wrap(errPanicInternal, fmt.Sprint(v))
			}

			c.Error(err)

			resp.Status = http.StatusInternalServerError

			e := m.tracer.Recovered(v)
			e.SetTransaction(tx)
			setContext(&e.Context, req, resp, body)
			e.Send()
		}

		if handlerErr != nil {
			e := m.tracer.NewError(handlerErr)
			setContext(&e.Context, req, resp, body)
			e.SetTransaction(tx)
			e.Handled = true
			e.Send()
		}

		tx.Result = apmhttp.StatusCodeResult(resp.Status)
		if tx.Sampled() {
			setContext(&tx.Context, req, resp, body)
		}

		body.Discard()
	}()

	handlerErr = m.handler(c)

	if handlerErr == nil {
		if !resp.Committed {
			resp.WriteHeader(http.StatusOK)
		}

		return nil
	}

	var errEchoHTTP *echo.HTTPError

	if ok := errors.As(handlerErr, &errEchoHTTP); ok {
		handlerErr = errEchoHTTP
		resp.Status = errEchoHTTP.Code

		reqPath := req.URL.RawPath
		if reqPath == "" {
			reqPath = req.URL.Path
		}

		if c.Path() != reqPath {
			return handlerErr
		}

		// When c.Path() matches the request path exactly,
		// that means either there's no matching route, or
		// there's an exactly matching route.
		//
		// When ErrNotFound or ErrMethodNotAllowed are
		// returned, it's probably because there's no
		// matching route, as opposed to the handler
		// returning them. We can confirm this by looking
		// for exact-matching routes.
		var unknownRoute bool

		if errors.Is(handlerErr, echo.ErrNotFound) {
			unknownRoute = isNotFoundHandler(c.Handler())
		}

		if errors.Is(handlerErr, echo.ErrMethodNotAllowed) {
			unknownRoute = isMethodNotAllowedHandler(c.Handler())
		}

		if unknownRoute {
			tx.Name = apmhttp.UnknownRouteRequestName(req)
		}
	}

	if c.Response().Committed {
		resp.Status = c.Response().Status
	}

	return handlerErr
}

func setContext(ctx *apm.Context, req *http.Request, resp *echo.Response, body *apm.BodyCapturer) {
	ctx.SetFramework("echo", echo.Version)
	ctx.SetHTTPRequest(req)
	ctx.SetHTTPRequestBody(body)
	ctx.SetHTTPStatusCode(resp.Status)
	ctx.SetHTTPResponseHeaders(resp.Header())
}

type options struct {
	tracer         *apm.Tracer
	requestIgnorer apmhttp.RequestIgnorerFunc
}

// APMOption sets options for tracing.
type APMOption func(*options)

// WithTracer returns an APMOption which sets t as the tracer
// to use for tracing server requests.
func WithTracer(t *apm.Tracer) APMOption {
	if t == nil {
		panic("t == nil")
	}

	return func(o *options) {
		o.tracer = t
	}
}

// WithRequestIgnorer returns a APMOption which sets r as the
// function to use to determine whether or not a request should
// be ignored. If r is nil, all requests will be reported.
func WithRequestIgnorer(r apmhttp.RequestIgnorerFunc) APMOption {
	if r == nil {
		r = apmhttp.IgnoreNone
	}

	return func(o *options) {
		o.requestIgnorer = r
	}
}

func isNotFoundHandler(h echo.HandlerFunc) bool {
	return isHandler(h, notFoundHandlerIdentity, &echo.NotFoundHandler)
}

func isMethodNotAllowedHandler(h echo.HandlerFunc) bool {
	return isHandler(h, methodNotAllowedHandlerIdentity, &echo.MethodNotAllowedHandler)
}

func isHandler(h echo.HandlerFunc, ident handlerFuncIdentity, handlerVar *func(echo.Context) error) bool {
	rv := reflect.ValueOf(h)

	ptr := rv.Pointer()
	if ptr == ident.rv.Pointer() {
		return true
	}

	// A sufficiently smart compiler could perform whole program optimisation
	// to determine that echo.NotFoundHandler and/or echo.MethodNotAllowedHandler
	// are only written to once to a defined function, enabling callers to inline
	// the assigned function. In this case, the function PC will not match.
	name := runtime.FuncForPC(ptr).Name()
	if name == ident.name {
		return true
	}

	// The global variables could have been reassigned since we read
	// their values during package init.
	ident = getHandlerFuncIdentity(*handlerVar)

	return ptr == ident.rv.Pointer() || name == ident.name
}

var (
	notFoundHandlerIdentity         = getHandlerFuncIdentity(echo.NotFoundHandler)
	methodNotAllowedHandlerIdentity = getHandlerFuncIdentity(echo.MethodNotAllowedHandler)
)

type handlerFuncIdentity struct {
	rv   reflect.Value
	name string
}

func getHandlerFuncIdentity(h func(echo.Context) error) handlerFuncIdentity {
	rv := reflect.ValueOf(h)

	return handlerFuncIdentity{
		rv:   rv,
		name: runtime.FuncForPC(rv.Pointer()).Name(),
	}
}
