package log

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

// NewLogger logger.
// If static fields are provided those values will define
// the default static fields for each new built instance
// if they were not yet configured.
func NewLogger(level Level, name string, fileLogger *lumberjack.Logger, batchCfg *BatchConfig, stfields ...interface{}) *Logger {
	if level < LevelDisabled || level > LevelError {
		level = LevelInfo
	}

	var (
		stdWriter io.Writer
		errWriter io.Writer
	)

	if fileLogger != nil {
		stdWriter = io.MultiWriter(os.Stdout, fileLogger)
		errWriter = io.MultiWriter(os.Stderr, fileLogger)
	} else {
		stdWriter = os.Stdout
		errWriter = os.Stderr
	}

	if batchCfg != nil {
		stdWriter = diode.NewWriter(stdWriter, batchCfg.MaxLines, batchCfg.Interval, func(missed int) {
			fmt.Printf("Logger Dropped %d messages", missed)
		})
		errWriter = diode.NewWriter(errWriter, batchCfg.MaxLines, batchCfg.Interval, func(missed int) {
			fmt.Printf("Logger Dropped %d messages", missed)
		})
	}

	stdl := log.Output(stdWriter).With().
		Timestamp().
		CallerWithSkipFrameCount(cfgSkipCallerCount).
		Stack().
		Logger()
	errl := log.Output(errWriter).With().
		Timestamp().
		CallerWithSkipFrameCount(cfgSkipCallerCount).
		Stack().
		Logger()

	setLogLevel(&stdl, level)
	setLogLevel(&errl, level)

	l := &Logger{
		Level:  level,
		StdLog: stdl,
		ErrLog: errl,
	}

	if len(stfields) > 1 && !cfg.configured {
		setup(level, name, fileLogger, false, batchCfg, stfields)

		defaultLogger = l
	}

	return l
}

// NewDevLogger logger.
// Pretty logging for development mode.
// Not recommended for production use.
// If static fields are provided those values will define
// the default static fields for each new built instance
// if they were not yet configured.
func NewDevLogger(level Level, name string, fileLogger *lumberjack.Logger, batchCfg *BatchConfig, stfields ...interface{}) *Logger {
	if level < LevelDisabled || level > LevelError {
		level = LevelInfo
	}

	var writer io.Writer

	if fileLogger != nil {
		writer = io.MultiWriter(os.Stdout, fileLogger)
	} else {
		writer = os.Stdout
	}

	if batchCfg != nil {
		writer = diode.NewWriter(writer, batchCfg.MaxLines, batchCfg.Interval, func(missed int) {
			fmt.Printf("Logger Dropped %d messages", missed)
		})
	}

	output := zerolog.ConsoleWriter{Out: writer, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("** %s **", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s=", i)
	}
	output.FormatErrFieldValue = func(i interface{}) string {
		if e, ok := i.(error); ok {
			return e.Error()
		}

		return fmt.Sprintf("%s", i)
	}
	output.FormatTimestamp = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	output.FormatErrFieldName = func(i interface{}) string {
		return "error="
	}
	output.FormatCaller = func(i interface{}) string {
		var c string

		if cc, ok := i.(string); ok {
			c = cc
		}

		if len(c) == 0 {
			return fmt.Sprintf("%s |", c)
		}

		cwd, err := os.Getwd()
		if err == nil {
			c = strings.TrimPrefix(c, cwd)
			c = strings.TrimPrefix(c, "/")
		}

		return fmt.Sprintf("%s |", c)
	}

	stdl := zerolog.New(output).With().
		Timestamp().
		CallerWithSkipFrameCount(cfgSkipCallerCount).
		Stack().
		Logger()
	errl := zerolog.New(output).With().
		Timestamp().
		CallerWithSkipFrameCount(cfgSkipCallerCount).
		Stack().
		Logger()

	setLogLevel(&stdl, level)
	setLogLevel(&errl, level)

	l := &Logger{
		Level:  level,
		StdLog: stdl,
		ErrLog: errl,
	}

	// if len(stfields) > 1 && !cfg.configured {
	if !cfg.configured {
		setup(level, name, fileLogger, true, batchCfg, stfields)

		defaultLogger = l
	}

	return l
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
