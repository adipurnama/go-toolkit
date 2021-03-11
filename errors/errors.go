package errors

import (
	"fmt"

	"github.com/adipurnama/go-toolkit/runtimekit"
	"github.com/iancoleman/strcase"
	pkg_errors "github.com/pkg/errors"
)

// 1. original caller function
// 2. WrapFunc / WrapFuncMsg function
// 3. wrapError().
const skipFnCall = 3

// New delegates to errors.New(...) from `github.com/pkg/errors`.
func New(s string) error {
	return pkg_errors.New(s)
}

// Is delegates to errors.Is(...) from `github.com/pkg/errors`.
func Is(err, target error) bool {
	return pkg_errors.Is(err, target)
}

// As delegates to errors.As(...) from `github.com/pkg/errors`.
func As(err error, target interface{}) bool {
	return pkg_errors.As(err, target)
}

// Cause delegates to errors.Cause(...) from `github.com/pkg/errors`.
func Cause(err error) error {
	return pkg_errors.Cause(err)
}

// WrapFunc wraps error with with file, line, & caller functions
//	with additional fields
// e.g. service.go; line 41; package: service; func Call();
//		errors.WrapFunc(errors.New("my error"), "param", "myValue")
//	=> service.go:41-service.Call param=myValue: my error
func WrapFunc(err error, keyVals ...interface{}) error {
	msg := runtimekit.CallerLineInfo(skipFnCall)
	return wrapError(err, msg, keyVals...)
}

// WrapFuncMsg wraps error with with file, line, & caller functions
//	with additional fields
// e.g. service.go; line 41; package: service; func Call();
//		errors.WrapFuncMsg(errors.New("my error"), "something happened", "param", "myValue")
//	=> service.go:41-service.Call something happened param=myValue: my error
func WrapFuncMsg(err error, msg string, keyVals ...interface{}) error {
	msg = fmt.Sprintf("%s %s", runtimekit.CallerLineInfo(skipFnCall), msg)

	return wrapError(err, msg, keyVals...)
}

func wrapError(err error, msg string, keyVals ...interface{}) error {
	if len(keyVals) <= 1 || (len(keyVals) > 1 && len(keyVals)%2 != 0) {
		return pkg_errors.Wrap(err, fmt.Sprintf("|> %s", msg))
	}

	for i := 0; i < len(keyVals)-1; i += 2 {
		if keyVals[i] == nil {
			continue
		}

		k := stringify(keyVals[i])
		v := keyVals[i+1]

		msg = fmt.Sprintf("%s %s=%v", msg, k, v)
	}

	return pkg_errors.Wrapf(err, fmt.Sprintf("=> %s", msg))
}

func stringify(val interface{}) string {
	switch v := val.(type) {
	case nil:
		return fmt.Sprintf("%v", v)
	case int:
		return fmt.Sprintf("%d", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case string:
		return strcase.ToSnake(v)
	default:
		return strcase.ToSnake(fmt.Sprintf("%+v", v))
	}
}
