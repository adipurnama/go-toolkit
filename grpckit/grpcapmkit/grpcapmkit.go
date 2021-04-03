package grpcapmkit

import (
	"context"

	"github.com/adipurnama/go-toolkit/runtimekit"
	"go.elastic.co/apm"
)

const skipFuncCount = 2

// ServerSpan retrieve span for grpc.Server handler.
func ServerSpan(ctx context.Context) *apm.Span {
	tx := apm.TransactionFromContext(ctx)
	return tx.StartSpan(runtimekit.SkippedFunctionName(skipFuncCount), "grpc.server", nil)
}
