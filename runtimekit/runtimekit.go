package runtimekit

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

const moduleName = "github.com/adipurnama/go-toolkit/"

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
		errFnLineInfo = strings.Replace(errFnLineInfo, moduleName, "", 1)
	}

	return errFnLineInfo
}

// CallerName returns this functions caller's name
// useful to wrap span, trace, context info
// e.g. trace.Start(ctx, runtimekit.CallerName()).
func CallerName() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}

	fn := runtime.FuncForPC(pc)
	result := fn.Name()

	cwd, err := os.Getwd()
	if err == nil {
		result = strings.TrimPrefix(result, cwd)
		result = strings.Replace(result, moduleName, "", 1)
	}

	return result
}
