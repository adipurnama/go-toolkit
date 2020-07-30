// Package log provide log interface to logging library
package log

import (
	"context"
	"fmt"
	stdLog "log"
	"os"
	"runtime"
	"strings"

	"github.com/pkg/errors"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	cfg config

	// defaultLogger is the package default logger.
	// It can be used right out of the box.
	// It can be replaced by a custom configured one
	// using package Set(*Logger) function
	// or using *Logger.Set() method.
	defaultLogger *Logger
)

func init() {
	defaultLogger = NewLogger(LevelDebug, "logger", nil)
}

// Logger is structured leveled logger
type Logger struct {
	// Level of min logging
	Level Level
	// Version
	Version string
	// Revision
	Revision string
	// DebugLog logger
	StdLog zerolog.Logger
	// ErrorLog logger
	ErrLog zerolog.Logger
	// Dynamic fields
	dynafields []interface{}
}

type config struct {
	// Name
	name string
	// Level of min logging
	level Level
	// Static fields
	stfields []interface{}
	// configured
	configured bool

	fileLogger *lumberjack.Logger

	isDevelopment bool
}

// ErrFunc is any function with empty argument and could returns error
type ErrFunc func() error

type contextKey struct {
	name string
}

// String returns formatted context key identifier.
func (k *contextKey) String() string {
	return "mw-" + k.name
}

// setup name and static fields.
// Each new instance of logger will always append these
// key-value pairs to the output and name if it is not empty.
// These values cannot be modified after they are configured.
func setup(name string, fileLogger *lumberjack.Logger, isDevelopment bool, stfields []interface{}) {
	if cfg.configured {
		return
	}

	cfg.isDevelopment = isDevelopment
	cfg.name = name
	cfg.stfields = append(cfg.stfields, stfields...)
	cfg.configured = true
	cfg.fileLogger = fileLogger
}

// SetFields set logger dynamic fields.
// The receiver instance will always append these
// key-value pairs to the output.
func (l *Logger) SetFields(dynafields ...interface{}) {
	l.dynafields = make([]interface{}, 0)
	l.dynafields = append(l.dynafields, dynafields...)
}

// AddField add dynamic field key-value
// The receiver instance will always append these
// key-value pairs to the output.
func (l *Logger) AddField(key, value interface{}) {
	l.dynafields = append(l.dynafields, key, value)
}

// ResetFields clear all the logger's  assigned dymanic fields
// Remove dynamic fields.
func (l *Logger) ResetFields() {
	l.dynafields = make([]interface{}, 0)
}

// GetLevelFromString return error level based on config string
func GetLevelFromString(level string) Level {
	switch level {
	case "info":
		return LevelInfo
	case "warn":
		return LevelWarn
	case "error":
		return LevelError
	default:
		return LevelDebug
	}
}

// AddToContext returns new context.Context with additional logger
func AddToContext(ctx context.Context, l *Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, l)
}

// Set the base package logger.
func Set(l *Logger) {
	defaultLogger = l
}

// Set default package logger.
// Can be used chained with NewLogger to create a new one,
// set it up as package default logger and get it for use in one step.
// i.e:
// logger := log.NewLogger(log.Debug, "name", "version", "revision").Set()
func (l *Logger) Set() *Logger {
	defaultLogger = l
	return defaultLogger
}

// NewContextLogger returns a copy of context that also includes a configured logger.
func NewContextLogger(ctx context.Context, fields ...interface{}) context.Context {
	l := FromCtx(ctx)

	if len(fields) > 0 {
		l.SetFields(fields...)
	}

	return context.WithValue(ctx, loggerCtxKey, l)
}

// FromCtx returns current logger in context.
// If there is not logger in context it returns
// a new one with current config values.
func FromCtx(ctx context.Context) *Logger {
	l, ok := ctx.Value(loggerCtxKey).(*Logger)
	if !ok {
		if cfg.isDevelopment {
			return NewDevLogger(cfg.level, cfg.name, cfg.fileLogger, cfg.stfields...)
		}

		return NewLogger(cfg.level, cfg.name, cfg.fileLogger, cfg.stfields...)
	}

	return l
}

// Debug logs debug messages.
func (l Logger) Debug(msg string, meta ...interface{}) {
	l.debugf(msg, meta)
}

// Info logs info messages.
func (l Logger) Info(msg string, meta ...interface{}) {
	l.infof(msg, meta)
}

// Warn logs warning messages.
func (l Logger) Warn(msg string, meta ...interface{}) {
	l.warnf(msg, meta)
}

// Error logs error messages.
func (l Logger) Error(err error, msg string, meta ...interface{}) {
	l.errorf(err, msg, meta)
}

func (l Logger) debugf(message string, fields []interface{}) {
	if l.Level > LevelDebug {
		return
	}

	le := l.StdLog.Debug()
	appendKeyValues(le, l.dynafields, fields)
	le.Msg(message)
}

func (l Logger) infof(message string, fields []interface{}) {
	if l.Level > LevelInfo {
		return
	}

	le := l.StdLog.Info()
	appendKeyValues(le, l.dynafields, fields)
	le.Msg(message)
}

func (l Logger) warnf(message string, fields []interface{}) {
	if l.Level > LevelWarn {
		return
	}

	le := l.StdLog.Warn()
	appendKeyValues(le, l.dynafields, fields)
	le.Msg(message)
}

func (l Logger) errorf(err error, message string, fields []interface{}) {
	le := l.ErrLog.Error()
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
			if cfg.stfields[i] != nil {
				k := stringify(cfg.stfields[i])
				fs[k] = cfg.stfields[i+1]
				i++
			}
		}
	}

	// check at least have "key":  "value"
	if len(dynafields) > 1 {
		for i := 0; i < len(dynafields)-1; i++ {
			if dynafields[i] != nil {
				k := stringify(dynafields[i])
				fs[k] = dynafields[i+1]
				// fmt.Printf("dyna - (%s, %v)\n", k, fs[k])
				i++
			}
		}
	}

	if len(fields) > 1 {
		for i := 0; i < len(fields)-1; i++ {
			if fields[i] != nil {
				k := stringify(fields[i])
				fs[k] = fields[i+1]
				// fmt.Printf("field - (%s, %v)\n", k, fs[k])
				i++
			}
		}
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
		return v
	default:
		return fmt.Sprintf("%+v", v)
	}
}

// UpdateLogLevel updates log level.
func (l Logger) UpdateLogLevel(level Level) {
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
	if fileFmt, ok := fileLineLogFmt("INFO"); ok {
		stdLog.Print(fileFmt + fmt.Sprint(v...))
		return
	}

	_ = stdLog.Output(cfgDefaultStdLogSkipCallerCount, fmt.Sprint(v...))
}

// Printf calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	if fileFmt, ok := fileLineLogFmt("INFO"); ok {
		stdLog.Printf(fileFmt+format, v...)
		return
	}

	_ = stdLog.Output(cfgDefaultStdLogSkipCallerCount, fmt.Sprintf(format+"\n", v...))
}

// Println calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Println.
func Println(v ...interface{}) {
	if fileFmt, ok := fileLineLogFmt("INFO"); ok {
		stdLog.Printf(fileFmt + fmt.Sprintln(v...))
		return
	}

	_ = stdLog.Output(cfgDefaultStdLogSkipCallerCount, fmt.Sprintln(v...))
}

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func Fatal(v ...interface{}) {
	if fileFmt, ok := fileLineLogFmt("ERROR"); ok {
		stdLog.Fatal(fileFmt + fmt.Sprintln(v...))
		return
	}

	_ = stdLog.Output(cfgDefaultStdLogSkipCallerCount, fmt.Sprintln(v...))

	os.Exit(1)
}

// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
func Fatalf(format string, v ...interface{}) {
	if fileFmt, ok := fileLineLogFmt("ERROR"); ok {
		stdLog.Fatalf(fileFmt+format+"\n", v...)
		return
	}

	_ = stdLog.Output(cfgDefaultStdLogSkipCallerCount, fmt.Sprintf(format+"\n", v...))

	os.Exit(1)
}

// helper function to get caller file & line number info
// returns formatted string ` | $LEVEL | dir/file:line | `
func fileLineLogFmt(level string) (string, bool) {
	skip := 3

	_, file, line, ok := runtime.Caller(skip)
	if ok {
		result := fmt.Sprintf("%s:%d", file, line)

		cwd, err := os.Getwd()
		if err == nil {
			result = strings.TrimPrefix(result, cwd)
			result = strings.TrimPrefix(result, "/")
		}

		return fmt.Sprintf(" | %s | %s | ", level, result), true
	}

	return "", false
}

// OnErrorFuncf execute function f with possible error return.
// It logs an error if error from f is found
// it aims to be used for deferred method such : resp.Body.Close(), tx.Rollback()
// Pass nil as logger to use default package logger
// So it should be used like: `defer log.OnOnErrorFunc(logger, resp.Body.Close, "closing response body for path %s", path)
func OnErrorFuncf(logger *Logger, f ErrFunc, format string, v ...interface{}) {
	err := f()
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

// OnErrorFunc execute function f with possible error return.
// It logs an error if error from f is found
// it aims to be used for deferred method such : resp.Body.Close(), tx.Rollback()
// Pass nil as logger to use default package logger
// So it should be used like: `defer log.OnOnErrorFunc(logger, resp.Body.Close)
func OnErrorFunc(logger *Logger, f ErrFunc) {
	err := f()
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
