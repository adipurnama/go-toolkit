package main

import (
	"context"

	"github.com/adipurnama/go-toolkit/example/logging/fakesql"
	"github.com/adipurnama/go-toolkit/example/logging/handler"
	"github.com/adipurnama/go-toolkit/example/logging/repository"
	"github.com/adipurnama/go-toolkit/example/logging/service"
	"github.com/adipurnama/go-toolkit/log"
	"github.com/adipurnama/go-toolkit/runtimekit"
	"github.com/pkg/errors"
)

const (
	isProdMode = false
)

var errPackageDomain = errors.New("original error")

func main() {
	if isProdMode {
		_ = log.NewLogger(
			log.LevelDebug,
			"sample-prod-logger",
			nil, nil,
			"custom", "value").Set()
	} else {
		_ = log.NewDevLogger(nil, nil, "example", true).Set()
	}

	ctx := context.Background()

	caller := runtimekit.FunctionName()
	caller2 := runtimekit.CallerLineInfo(1)

	err1 := errors.Wrapf(errPackageDomain, "additional message=%s val=%s key=%s", "val0", "val", "val2")
	log.Println("caller: ", caller, "caller2: ", caller2, "error1=", err1)

	logger := log.FromCtx(ctx)
	logger.AddField("my_field", "custom-field-value")
	ctx = log.AddToContext(ctx, logger)

	log.FromCtx(ctx).Warn("log warn message", "Key-Random1", "value")
	log.FromCtx(ctx).Info("log info message - no error", "field_key1", "whatever")

	// 3-tier architecture error tracing
	r := &repository.MockDBRepository{DB: &fakesql.DB{}}
	s := &service.Service{R: r}
	h := &handler.Handler{S: s}
	validUserID := 10
	invalidUserID := 69

	err := h.FindUserByID(validUserID)
	log.FromCtx(ctx).Error(err, "log debug message", "field_key1", "whatever")

	err = h.FindUserByID(invalidUserID)
	log.FromCtx(ctx).Error(err, "log debug message", "field_key1", "whatever")

	log.FromCtx(ctx).Info("log info message - no error", "field_key1", "whatever")
}
