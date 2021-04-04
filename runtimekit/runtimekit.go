package runtimekit

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
)

// NewRuntimeContext returns context & cancel func listening to :
// - os.Interrupt
// - syscall.SIGTERM
// - syscall.SIGINT.
func NewRuntimeContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
}

// CallerLineInfo returns caller file:line-package.function
// e.g. service.go:38-service.CallAPI.
func CallerLineInfo(skip int) string {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}

	files := strings.Split(file, "/")
	file = files[len(files)-1]

	fn := runtime.FuncForPC(pc)
	fnNames := strings.Split(fn.Name(), "/")
	fnName := fnNames[len(fnNames)-1]

	errFnLineInfo := fmt.Sprintf("%s:%d-%s", file, line, fnName)

	cwd, errGetWd := os.Getwd()
	if errGetWd == nil {
		errFnLineInfo = strings.TrimPrefix(errFnLineInfo, cwd)
		errFnLineInfo = strings.TrimPrefix(errFnLineInfo, "/")
		moduleFnNames := strings.Split(errFnLineInfo, "/")
		errFnLineInfo = moduleFnNames[len(moduleFnNames)-1]
	}

	return errFnLineInfo
}

// FunctionName returns this function caller's name
// useful to wrap span, trace, context info
// e.g. trace.Start(ctx, runtimekit.FunctionName()).
func FunctionName() string {
	skipCount := 2
	return SkippedFunctionName(skipCount)
}

// SkippedFunctionName returns function caller's name with skipped count.
func SkippedFunctionName(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}

	fn := runtime.FuncForPC(pc)
	result := fn.Name()

	cwd, err := os.Getwd()
	if err == nil {
		result = strings.TrimPrefix(result, cwd)
		moduleFnNames := strings.Split(result, "/")
		result = moduleFnNames[len(moduleFnNames)-1]
	}

	return result
}
