package log

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

// NewLogger logger.
// If static fields are provided those values will define
// the default static fields for each new built instance
// if they were not yet configured.
func NewLogger(level int, name string, fileLogger *lumberjack.Logger, stfields ...interface{}) *Logger {
	if level < Disabled || level > LevelError {
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

	stdl := log.Output(stdWriter).With().
		Timestamp().
		CallerWithSkipFrameCount(4).
		Stack().
		Logger()
	errl := log.Output(errWriter).With().
		Timestamp().
		CallerWithSkipFrameCount(4).
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
		// if !cfg.configured {
		setup(name, fileLogger, false, stfields)

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
func NewDevLogger(level int, name string, fileLogger *lumberjack.Logger, stfields ...interface{}) *Logger {
	if level < Disabled || level > LevelError {
		level = LevelInfo
	}

	var output zerolog.ConsoleWriter

	if fileLogger != nil {
		multiWriter := io.MultiWriter(os.Stdout, fileLogger)
		output = zerolog.ConsoleWriter{Out: multiWriter, TimeFormat: time.RFC3339}
	} else {
		output = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	}

	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s |", i)
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
	output.FormatErrFieldName = func(i interface{}) string {
		return "error="
	}
	output.FormatCaller = func(i interface{}) string {
		var c string

		if cc, ok := i.(string); ok {
			c = cc
		}

		if len(c) > 0 {
			cwd, err := os.Getwd()
			if err == nil {
				c = strings.TrimPrefix(c, cwd)
				c = strings.TrimPrefix(c, "/")
			}
		}

		return fmt.Sprintf("%s |", c)
	}

	stdl := zerolog.New(output).With().
		Timestamp().
		CallerWithSkipFrameCount(4).
		Stack().
		Logger()
	errl := zerolog.New(output).With().
		Timestamp().
		CallerWithSkipFrameCount(4).
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
		setup(name, fileLogger, true, stfields)

		defaultLogger = l
	}

	return l
}
