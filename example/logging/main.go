package main

import (
	"context"

	"github.com/adipurnama/go-toolkit/log"
)

func main() {
	_ = log.NewDevLogger(log.LevelDebug, "sample-logger", nil, "example", true).Set()

	ctx := context.Background()
	logger := log.FromCtx(ctx)
	logger.AddField("my_field", "custom")
	ctx = log.NewContextLogger(ctx, logger)

	log.FromCtx(ctx).Info("debug message - no error", "field_here", "whatever")

	log.FromCtx(ctx).Error(definitelyError(), "debug message", "field_here", "whatever")

	log.FromCtx(ctx).Info("debug message - no error", "field_here", "whatever")
}
