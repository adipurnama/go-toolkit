// Package tracer provides apm tracing helper
package tracer

import (
	"context"

	"github.com/adipurnama/go-toolkit/runtimekit"
	"go.elastic.co/apm"
)

const skipFuncCount = 2

// ServiceFuncSpan returns span for service layer type.
func ServiceFuncSpan(ctx context.Context) *apm.Span {
	tx := apm.TransactionFromContext(ctx)
	return tx.StartSpan(runtimekit.SkippedFunctionName(skipFuncCount), "service", nil)
}

// RepositoryFuncSpan returns span for repository layer type.
func RepositoryFuncSpan(ctx context.Context) *apm.Span {
	tx := apm.TransactionFromContext(ctx)
	return tx.StartSpan(runtimekit.SkippedFunctionName(skipFuncCount), "repository", nil)
}

// APIClientFuncSpan returns span for http/grpc.Client layer type.
func APIClientFuncSpan(ctx context.Context) *apm.Span {
	tx := apm.TransactionFromContext(ctx)
	return tx.StartSpan(runtimekit.SkippedFunctionName(skipFuncCount), "api_client", nil)
}
