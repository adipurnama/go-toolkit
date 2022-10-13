// Package log is internal log wrapper functionality
package log

import (
	"context"
	"fmt"
	"io"
	stdLog "log"
	"os"
	"runtime"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"

	"github.com/rs/zerolog"
)

// defaultLogger is the package default logger.
// It can be used right out of the box.
// It can be replaced by a custom configured one
// using package Set(*Logger) function
// or using *Logger.Set() method.
var defaultLogger *Logger

func init() {
	stdLog.SetOutput(os.Stdout)

	zerolog.ErrorStackMarshaler = marshalStack
	zerolog.ErrorStackFieldName = "stacktrace"

	defaultLogger = NewLogger(LevelDebug, "logger", nil, nil)
}

// AddToContext returns new context.Context with additional logger.
func AddToContext(ctx context.Context, l *Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, l)
}

// NewLoggingContext returns a copy of context that also includes a configured logger.
func NewLoggingContext(ctx context.Context, fields ...interface{}) context.Context {
	l := FromCtx(ctx)

	if len(fields) > 0 {
		l.SetFields(fields...)
	}

	return context.WithValue(ctx, loggerCtxKey, l)
}

// FromCtx returns current logger in context.
// If there is no logger in context it returns
// a new one with current config values.
// logger initial attribute fields is copied from existing defaultLogger fields.
func FromCtx(ctx context.Context) *Logger {
	if l, ok := ctx.Value(loggerCtxKey).(*Logger); ok {
		return l
	}

	fields := make([]interface{}, len(defaultLogger.dynafields))
	_ = copy(fields, defaultLogger.dynafields)

	return &Logger{
		Level:      defaultLogger.Level,
		Version:    defaultLogger.Revision,
		Revision:   defaultLogger.Revision,
		StdLog:     defaultLogger.StdLog,
		ErrLog:     defaultLogger.ErrLog,
		dynafields: fields,
	}
}

// Debug logs debug messages.
func (l *Logger) Debug(msg string, meta ...interface{}) {
	l.debugf(msg, meta)
}

// Info logs info messages.
func (l *Logger) Info(msg string, meta ...interface{}) {
	l.infof(msg, meta)
}

// Warn logs warning messages.
func (l *Logger) Warn(msg string, meta ...interface{}) {
	l.warnf(nil, msg, meta)
}

// Error logs error messages.
func (l *Logger) Error(err error, msg string, meta ...interface{}) {
	if defaultLogger.logFmt {
		stdLog.Println("-----------------------")
	}

	l.errorf(err, msg, meta)

	if defaultLogger.logFmt {
		stdLog.Println("-----------------------")
	}
}

// WarnError used for log error but in `warn` level
// e.g.
//
//	if err != nil {
//	  log.FromCtx(ctx).WarnError(err, "something happened. continue...")
//	}.
func (l *Logger) WarnError(err error, msg string, meta ...interface{}) {
	if defaultLogger.logFmt {
		stdLog.Println("-----------------------")
	}

	l.warnf(err, msg, meta)

	if defaultLogger.logFmt {
		stdLog.Println("-----------------------")
	}
}

func (l *Logger) debugf(message string, fields []interface{}) {
	if l.Level > LevelDebug {
		return
	}

	le := l.StdLog.Debug().Stack()
	appendKeyValues(le, l.dynafields, fields)
	le.Msg(message)
}

func (l *Logger) infof(message string, fields []interface{}) {
	if l.Level > LevelInfo {
		return
	}

	le := l.StdLog.Info().Stack()
	appendKeyValues(le, l.dynafields, fields)
	le.Msg(message)
}

func (l *Logger) warnf(err error, message string, fields []interface{}) {
	if l.Level > LevelWarn {
		return
	}

	le := l.StdLog.Warn().Stack()
	appendKeyValues(le, l.dynafields, fields)

	if err != nil {
		le.Err(err)
	}

	le.Msg(message)
}

func (l *Logger) errorf(err error, message string, fields []interface{}) {
	le := l.ErrLog.Error().Stack()
	appendKeyValues(le, l.dynafields, fields)
	le.Err(err)
	le.Msg(message)
}

// TODO: Optimize.
// Static key-value calculation shoud be cached.
// Dynamic key-value calculation shoud be cached if didn't changed.
func appendKeyValues(le *zerolog.Event, dynafields []interface{}, fields []interface{}) {
	if cfg.name != "" {
		le.Str("name", cfg.name)
	}

	fs := make(map[string]interface{})

	if len(cfg.stfields) > 1 {
		for i := 0; i < len(cfg.stfields)-1; i++ {
			if cfg.stfields[i] == nil {
				continue
			}

			k := stringify(cfg.stfields[i])
			if IsSensitiveParam(k) {
				fs[k] = RedactionString
				i++

				continue
			}

			fs[k] = cfg.stfields[i+1]
			i++
		}
	}

	// check at least have "key":  "value"
	if len(dynafields) > 1 {
		for i := 0; i < len(dynafields)-1; i++ {
			if dynafields[i] == nil {
				continue
			}

			k := stringify(dynafields[i])
			if IsSensitiveParam(k) {
				fs[k] = RedactionString
				i++

				continue
			}

			fs[k] = dynafields[i+1]
			i++
		}
	}

	// check at least have 1 key:value
	if len(fields) <= 1 {
		le.Fields(fs)
		return
	}

	for i := 0; i < len(fields)-1; i++ {
		if fields[i] == nil {
			continue
		}

		k := stringify(fields[i])
		if IsSensitiveParam(k) {
			fs[k] = RedactionString
			i++

			continue
		}

		if k == "error" {
			if errVal, ok := fields[i+1].(error); ok {
				le.Err(errVal)
				i++

				continue
			}
		}

		fs[k] = fields[i+1]
		i++
	}

	le.Fields(fs)
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

// UpdateLogLevel updates log level.
func (l *Logger) UpdateLogLevel(level Level) {
	// Allow info level to log the update
	// But don't downgrade to it if Error is set.
	current := LevelError

	l.Info("Log level updated", "", "log level", level)

	l.Level = current
	if level < LevelDisabled || level > LevelError {
		l.Level = level
		setLogLevel(&l.StdLog, level)
		setLogLevel(&l.ErrLog, level)
	}
}

func setLogLevel(l *zerolog.Logger, level Level) {
	switch level {
	case LevelDisabled:
		l.Level(zerolog.Disabled)
	case LevelDebug:
		l.Level(zerolog.DebugLevel)
	case LevelInfo:
		l.Level(zerolog.InfoLevel)
	case LevelWarn:
		l.Level(zerolog.WarnLevel)
	case LevelError:
		l.Level(zerolog.ErrorLevel)
	default:
		l.Level(zerolog.DebugLevel)
	}
}

// These functions write to the standard logger.

// Print calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) {
	if fileFmt, ok := fileLineLogFmt(LevelInfo); ok {
		stdLog.Print(fileFmt + fmt.Sprint(v...))
		return
	}

	_ = stdLog.Output(cfgDefaultStdLogSkipCallerCount, fmt.Sprint(v...))
}

// Printf calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	if fileFmt, ok := stdFileLineLogFmt(LevelInfo); ok {
		stdLog.Printf(fileFmt+format, v...)
		return
	}

	_ = stdLog.Output(cfgDefaultStdLogSkipCallerCount, fmt.Sprintf(format+"\n", v...))
}

// Println calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Println.
func Println(v ...interface{}) {
	if fileFmt, ok := stdFileLineLogFmt(LevelInfo); ok {
		stdLog.Printf(fileFmt + fmt.Sprintln(v...))
		return
	}

	_ = stdLog.Output(cfgDefaultStdLogSkipCallerCount, fmt.Sprintln(v...))
}

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func Fatal(v ...interface{}) {
	if fileFmt, ok := stdFileLineLogFmt(LevelError); ok {
		stdLog.Fatal(fileFmt + fmt.Sprintln(v...))
		return
	}

	_ = stdLog.Output(cfgDefaultStdLogSkipCallerCount, fmt.Sprintln(v...))

	os.Exit(1)
}

// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
func Fatalf(format string, v ...interface{}) {
	if fileFmt, ok := stdFileLineLogFmt(LevelError); ok {
		stdLog.Fatalf(fileFmt+format+"\n", v...)
		return
	}

	_ = stdLog.Output(cfgDefaultStdLogSkipCallerCount, fmt.Sprintf(format+"\n", v...))

	os.Exit(1)
}

// helper function to get caller file & line number info
// returns formatted string ` | $LEVEL | dir/file:line | `.
func fileLineLogFmt(level Level) (string, bool) {
	skipCallerCount := 3

	_, file, line, ok := runtime.Caller(skipCallerCount)
	if !ok {
		return "", false
	}

	files := strings.Split(file, "/")
	file = files[len(files)-1]

	result := fmt.Sprintf("%s:%d", file, line)

	cwd, err := os.Getwd()
	if err == nil {
		result = strings.TrimPrefix(result, cwd)
		result = strings.TrimPrefix(result, "/")
	}

	return fmt.Sprintf(" | %s | %s | ", levelString(level), result), true
}

func stdFileLineLogFmt(level Level) (string, bool) {
	skipCallerCount := 2

	_, file, line, ok := runtime.Caller(skipCallerCount)
	if !ok {
		return "", false
	}

	files := strings.Split(file, "/")
	file = files[len(files)-1]

	result := fmt.Sprintf("%s:%d", file, line)

	cwd, err := os.Getwd()
	if err == nil {
		result = strings.TrimPrefix(result, cwd)
		result = strings.TrimPrefix(result, "/")
	}

	return fmt.Sprintf(" | %s | %s | ", levelString(level), result), true
}

// OnCloseErrorf execute function f with possible error return.
// It logs an error if error from f is found
// it aims to be used for deferred method such : resp.Body.Close(), tx.Rollback()
// Pass nil as logger to use default package logger
// So it should be used like: `defer log.OnOnErrorFunc(logger, resp.Body.Close, "closing response body for path %s", path).
func OnCloseErrorf(logger *Logger, f io.Closer, format string, v ...interface{}) {
	err := f.Close()
	if err == nil {
		return
	}

	if errors.Is(err, os.ErrClosed) {
		return
	}

	if logger == nil {
		logger = defaultLogger
	}

	logger.Warn("error detected", "err", errors.Wrapf(err, format, v...))
}

// OnCloseError execute function f with possible error return.
// It logs an error if error from f is found
// it aims to be used for deferred method such : resp.Body.Close(), tx.Rollback()
// Pass nil as logger to use default package logger
// So it should be used like: `defer log.OnOnErrorFunc(logger, resp.Body.Close).
func OnCloseError(logger *Logger, f io.Closer) {
	err := f.Close()
	if err == nil {
		return
	}

	if errors.Is(err, os.ErrClosed) {
		return
	}

	if logger == nil {
		logger = defaultLogger
	}

	logger.Warn("error detected", "err", err)
}
