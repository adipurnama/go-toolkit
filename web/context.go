package web

import (
	"context"
	"net/http"
	"os"
	"syscall"

	"github.com/adipurnama/go-toolkit/log"
	"github.com/sethvargo/go-signalcontext"
)

var (
	// ContextKeyRequestID to obtains requestID from context
	ContextKeyRequestID = ContextKey("reqID")

	// ContextKeyHeader to obtains original http header from downstream
	ContextKeyHeader = ContextKey("header")

	// ContextKeyTraceID to obtains traceID from downstream context
	ContextKeyTraceID = ContextKey("traceID")
)

// NewRuntimeContext returns context & cancel func listening to :
// - os.Interrupt
// - syscall.SIGTERM
// - syscall.SIGINT
func NewRuntimeContext() (context.Context, context.CancelFunc) {
	return signalcontext.Wrap(
		log.NewLoggingContext(context.Background()),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
	)
}

// ContextID is a struct which will be used as context key.
type ContextID struct {
	name string
}

// String returns formatted context key identifier.
func (k *ContextID) String() string {
	return "web context: " + k.name
}

// ContextKey constructs context key using name supplied.
func ContextKey(name string) *ContextID {
	return &ContextID{name: name}
}

// ValueFromContext returns string value for certain ContextID.
func ValueFromContext(ctx context.Context, key *ContextID) string {
	if val, ok := ctx.Value(key).(string); ok {
		return val
	}

	return ""
}

// HeaderFromContext - get header value from context
// we set "header" key on context to set forwarded request context.
func HeaderFromContext(ctx context.Context) http.Header {
	if val, ok := ctx.Value(ContextKeyHeader).(http.Header); ok {
		return val
	}

	val := make(http.Header)

	return val
}
