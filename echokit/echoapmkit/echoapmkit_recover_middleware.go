package echoapmkit

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"runtime/debug"

	"github.com/adipurnama/go-toolkit/log"
	echo "github.com/labstack/echo/v4"
	"github.com/pinpoint-apm/pinpoint-go-agent"
	phttp "github.com/pinpoint-apm/pinpoint-go-agent/plugin/http"
	"github.com/pkg/errors"
	"go.elastic.co/apm"
	apmhttp "go.elastic.co/apm/module/apmhttp"
)

// MODIFIED from https://github.com/elastic/apm-agent-go/blob/master/module/apmechov4/middleware.go

// RecoverMiddleware returns a new Echo middleware handler for tracing
// requests and reporting errors.
//
// This middleware will recover and report panics, so it can
// be used instead of echo/middleware.Recover.
//
// By default, the middleware will use apm.DefaultTracer.
// Use WithTracer to specify an alternative tracer.
func RecoverMiddleware(o ...APMOption) echo.MiddlewareFunc {
	opts := options{
		elasticTracer:         apm.DefaultTracer,
		elasticRequestIgnorer: apmhttp.DefaultServerRequestIgnorer(),
	}

	for _, o := range o {
		o(&opts)
	}

	return func(h echo.HandlerFunc) echo.HandlerFunc {
		m := &middleware{
			elasticTracer:         opts.elasticTracer,
			pinpointAgent:         opts.pinpointAgent,
			handler:               h,
			elasticRequestIgnorer: opts.elasticRequestIgnorer,
		}

		return m.handle
	}
}

var errPanicInternal = errors.New("found panic while serving request")

type middleware struct {
	handler               echo.HandlerFunc
	elasticTracer         *apm.Tracer
	elasticRequestIgnorer apmhttp.RequestIgnorerFunc
	pinpointAgent         pinpoint.Agent
}

func (m *middleware) handle(c echo.Context) error {
	req := c.Request()
	name := req.Method + " " + c.Path()

	var (
		tx    *apm.Transaction
		eBody *apm.BodyCapturer
	)

	if m.elasticTracer.Recording() && !m.elasticRequestIgnorer(req) {
		tx, req = apmhttp.StartTransaction(m.elasticTracer, name, req)
		defer tx.End()

		c.SetRequest(req)

		eBody = m.elasticTracer.CaptureHTTPRequestBody(req)
	}

	var pTracer pinpoint.Tracer

	if m.pinpointAgent != nil && m.pinpointAgent.Enable() {
		pTracer = phttp.NewHttpServerTracer(m.pinpointAgent, req, "Echo Server")
		defer pTracer.EndSpan()

		ctx := pinpoint.NewContext(req.Context(), pTracer)

		c.SetRequest(req.WithContext(ctx))
		defer pTracer.NewSpanEvent(req.Method + " " + c.Path()).EndSpan()
	}

	resp := c.Response()

	resp.Status = http.StatusInternalServerError

	var errHandler error

	defer func() {
		if v := recover(); v != nil {
			errStack := debug.Stack()

			err, ok := v.(error)
			if !ok {
				err = errors.Wrap(errPanicInternal, fmt.Sprint(v))
			}

			log.FromCtx(c.Request().Context()).Error(
				err,
				"recovered from panic",
				"panic_stack", errStack,
			)

			c.Error(err)

			resp.Status = http.StatusInternalServerError

			if eBody != nil {
				e := m.elasticTracer.Recovered(v)
				e.SetTransaction(tx)
				setContext(&e.Context, req, resp, eBody)
				e.Send()
			}

			if pTracer != nil {
				pTracer.Span().SetError(err)
			}
		}

		if errHandler != nil {
			if eBody != nil {
				e := m.elasticTracer.NewError(errHandler)
				setContext(&e.Context, req, resp, eBody)
				e.SetTransaction(tx)
				e.Handled = true
				e.Send()
			}

			if pTracer != nil {
				pTracer.Span().SetError(errHandler)
			}
		}

		if eBody != nil {
			tx.Result = apmhttp.StatusCodeResult(resp.Status)
			if tx.Sampled() {
				setContext(&tx.Context, req, resp, eBody)
			}

			eBody.Discard()
		}

		if pTracer != nil {
			phttp.RecordHttpServerResponse(pTracer, c.Response().Status, c.Response().Header())
		}
	}()

	errHandler = m.handler(c)

	if errHandler == nil {
		if !resp.Committed {
			resp.WriteHeader(http.StatusOK)
		}

		return nil
	}

	var errEchoHTTP *echo.HTTPError

	if ok := errors.As(errHandler, &errEchoHTTP); ok {
		errHandler = errEchoHTTP
		resp.Status = errEchoHTTP.Code

		reqPath := req.URL.RawPath
		if reqPath == "" {
			reqPath = req.URL.Path
		}

		if c.Path() != reqPath {
			return errHandler
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

		if errors.Is(errHandler, echo.ErrNotFound) {
			unknownRoute = isNotFoundHandler(c.Handler())
		}

		if errors.Is(errHandler, echo.ErrMethodNotAllowed) {
			unknownRoute = isMethodNotAllowedHandler(c.Handler())
		}

		if unknownRoute {
			tx.Name = apmhttp.UnknownRouteRequestName(req)
		}
	}

	if c.Response().Committed {
		resp.Status = c.Response().Status
	}

	return errHandler
}

func setContext(ctx *apm.Context, req *http.Request, resp *echo.Response, body *apm.BodyCapturer) {
	ctx.SetFramework("echo", echo.Version)
	ctx.SetHTTPRequest(req)
	ctx.SetHTTPRequestBody(body)
	ctx.SetHTTPStatusCode(resp.Status)
	ctx.SetHTTPResponseHeaders(resp.Header())
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
