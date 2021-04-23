package web

import (
	"context"
	"net/http"

	ut "github.com/go-playground/universal-translator"
)

var (
	// ContextKeyRequestID to store/obtains requestID from context.
	ContextKeyRequestID = ContextKey("reqID")

	// ContextKeyHeader to store/obtains original http header from downstream.
	ContextKeyHeader = ContextKey("header")

	// ContextKeyTraceID to store/obtains traceID from downstream context.
	ContextKeyTraceID = ContextKey("traceID")

	// ContextKeyTranslator to store/obtains translator to/from request's context.
	ContextKeyTranslator = ContextKey("translator")
)

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

func translatorFromContext(ctx context.Context) ut.Translator {
	if val, ok := ctx.Value(ContextKeyTranslator).(ut.Translator); ok {
		return val
	}

	return nil
}
