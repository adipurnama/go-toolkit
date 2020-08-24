package main

import (
	"context"

	"github.com/adipurnama/go-toolkit/example/logging/helper"
	"github.com/adipurnama/go-toolkit/log"
)

func main() {
	// _ = log.NewLogger(log.LevelDebug, "sample-prod-logger", nil, "custom", "value").Set()
	_ = log.NewDevLogger(log.LevelDebug, "sample-logger", nil, nil, "example", true).Set()

	ctx := context.Background()
	logger := log.FromCtx(ctx)
	logger.AddField("my_field", "custom")
	ctx = log.NewContextLogger(ctx, logger)

	log.FromCtx(ctx).Info("debug message - no error", "field_here", "whatever")

	log.FromCtx(ctx).Error(helper.DefinitelyError(), "debug message", "field_here", "whatever")

	log.FromCtx(ctx).Info("debug message - no error", "field_here", "whatever")
}
