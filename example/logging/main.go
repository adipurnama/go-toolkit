package main

import (
	"context"

	"github.com/adipurnama/go-toolkit/errors"
	"github.com/adipurnama/go-toolkit/example/logging/helper"
	"github.com/adipurnama/go-toolkit/log"
	"github.com/adipurnama/go-toolkit/runtimekit"
)

const (
	isProdMode = false
)

var errOriginal = errors.New("original error")

func main() {
	if isProdMode {
		_ = log.NewLogger(
			log.LevelDebug,
			"sample-prod-logger",
			nil, nil,
			"custom", "value").Set()
	} else {
		_ = log.NewDevLogger(
			log.LevelDebug,
			"sample-logger",
			nil, nil,
			"example", true).Set()
	}

	ctx := context.Background()

	caller := runtimekit.CallerName()
	caller2 := runtimekit.CallerLineInfo(1)

	err1 := errors.WrapFunc(errOriginal)
	log.Println("errors.WrapFunc:", err1)

	err1 = errors.WrapFunc(err1, "token", "value")
	log.Println("errors.WrapFunc:", err1)

	err1 = errors.WrapFuncMsg(err1, "additional message", nil, "val", "key2", "val2")
	log.Println("caller: ", caller, "caller2: ", caller2, "error=", err1)

	logger := log.FromCtx(ctx)
	logger.AddField("my_field", "custom-field-value")
	ctx = log.AddToContext(ctx, logger)

	log.FromCtx(ctx).Warn("log warn message", "Key-Random1", "value")
	log.FromCtx(ctx).Info("log info message - no error", "field_key1", "whatever")

	log.FromCtx(ctx).Error(helper.DefinitelyError(), "log debug message", "field_key1", "whatever")

	log.FromCtx(ctx).Info("log info message - no error", "field_key1", "whatever")
}
