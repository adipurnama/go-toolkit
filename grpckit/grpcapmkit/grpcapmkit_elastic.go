package grpcapmkit

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"strings"

	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

var (
	elasticTraceparentHeader = strings.ToLower(apmhttp.ElasticTraceparentHeader)
	w3cTraceparentHeader     = strings.ToLower(apmhttp.W3CTraceparentHeader)
	tracestateHeader         = strings.ToLower(apmhttp.TracestateHeader)
)

func startTransaction(ctx context.Context, tracer *apm.Tracer, name string) (*apm.Transaction, context.Context) {
	var opts apm.TransactionOptions

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		traceContext, ok := getIncomingMetadataTraceContext(md, elasticTraceparentHeader)
		if !ok {
			traceContext, _ = getIncomingMetadataTraceContext(md, w3cTraceparentHeader)
		}

		opts.TraceContext = traceContext
	}

	tx := tracer.StartTransactionOptions(name, "request", opts)

	tx.Context.SetFramework("grpc", grpc.Version)

	if peer, ok := peer.FromContext(ctx); ok {
		// Set underlying HTTP/2.0 request context.
		//
		// https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md
		var (
			tlsConnectionState *tls.ConnectionState
			peerAddr           string
			authority          string
		)

		url := url.URL{Scheme: "http", Path: name}

		if info, ok := peer.AuthInfo.(credentials.TLSInfo); ok {
			url.Scheme = "https"
			tlsConnectionState = &info.State
		}

		if peer.Addr != nil {
			peerAddr = peer.Addr.String()
		}

		if values := md.Get(":authority"); len(values) > 0 {
			authority = values[0]
		}

		protoMajor := 2

		tx.Context.SetHTTPRequest(&http.Request{
			URL:        &url,
			Method:     "POST", // method is always POST
			ProtoMajor: protoMajor,
			ProtoMinor: 0,
			Header:     http.Header(md),
			Host:       authority,
			RemoteAddr: peerAddr,
			TLS:        tlsConnectionState,
		})
	}

	return tx, apm.ContextWithTransaction(ctx, tx)
}

func getIncomingMetadataTraceContext(md metadata.MD, header string) (apm.TraceContext, bool) {
	if values := md.Get(header); len(values) == 1 {
		traceContext, err := apmhttp.ParseTraceparentHeader(values[0])

		if err == nil {
			traceContext.State, _ = apmhttp.ParseTracestateHeader(md.Get(tracestateHeader)...)
			return traceContext, true
		}
	}

	return apm.TraceContext{}, false
}

func setTransactionResult(tx *apm.Transaction, err error) {
	statusCode := statusCodeFromError(err)
	tx.Result = statusCode.String()

	// For gRPC servers, the transaction outcome is generally "success",
	// except for codes which are not subject to client interpretation.
	if tx.Outcome == "" {
		switch statusCode {
		case codes.Unknown,
			codes.DeadlineExceeded,
			codes.ResourceExhausted,
			codes.FailedPrecondition,
			codes.Aborted,
			codes.Internal,
			codes.Unavailable,
			codes.DataLoss:
			tx.Outcome = "failure"
		default:
			tx.Outcome = "success"
		}
	}
}
