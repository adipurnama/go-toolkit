package log

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

// NewDevLogger logger.
// Pretty logging for development mode.
// Not recommended for production use.
// If static fields are provided those values will define
// the default static fields for each new built instance
// if they were not yet configured.
func NewDevLogger(fileLogger *lumberjack.Logger, batchCfg *BatchConfig, stfields ...interface{}) *Logger {
	var writer io.Writer

	if fileLogger != nil {
		writer = io.MultiWriter(os.Stdout, fileLogger)
	} else {
		writer = os.Stdout
	}

	if batchCfg != nil && batchCfg.Enabled {
		writer = diode.NewWriter(writer, batchCfg.MaxLines, batchCfg.Interval, func(missed int) {
			fmt.Printf("Logger Dropped %d messages", missed)
		})
	}

	output := zerolog.ConsoleWriter{Out: writer, TimeFormat: time.RFC3339}
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

	stdl := zerolog.New(output).With().
		Timestamp().
		CallerWithSkipFrameCount(cfgSkipCallerCount).
		Logger()
	errl := zerolog.New(output).With().
		Timestamp().
		CallerWithSkipFrameCount(cfgSkipCallerCount).
		Logger()

	level := LevelDebug

	setLogLevel(&stdl, level)
	setLogLevel(&errl, level)

	l := &Logger{
		Level:  level,
		StdLog: stdl,
		ErrLog: errl,
		logFmt: true,
	}

	// if len(stfields) > 1 && !cfg.configured {
	if !cfg.configured {
		setup(level, "", fileLogger, true, batchCfg, stfields)

		defaultLogger = l
	}

	return l
}

// NewLogger logger.
// If static fields are provided those values will define
// the default static fields for each new built instance
// if they were not yet configured.
func NewLogger(level Level, name string, fileLogger *lumberjack.Logger, batchCfg *BatchConfig, stfields ...interface{}) *Logger {
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

	if batchCfg != nil && batchCfg.Enabled {
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
		Logger()
	errl := log.Output(errWriter).With().
		Timestamp().
		CallerWithSkipFrameCount(cfgSkipCallerCount).
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

// Set the base package logger.
func Set(l *Logger) {
	defaultLogger = l
}

// Set default package logger.
// Can be used chained with NewLogger to create a new one,
// set it up as package default logger and get it for use in one step.
// i.e:
// logger := log.NewLogger(log.Debug, "name", "version", "revision").Set().
func (l *Logger) Set() *Logger {
	defaultLogger = l
	return defaultLogger
}
