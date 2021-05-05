package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/log"

	kitconfig "github.com/adipurnama/go-toolkit/config"
	"github.com/pkg/errors"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
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

	output.FormatCaller = func(i interface{}) string {
		var sb strings.Builder

		callers := strings.Split(fmt.Sprintf("%s", i), "/")

		for i, v := range callers {
			if i == len(callers)-1 {
				sb.WriteString(callers[i])
				return sb.String()
			}

			if i == len(callers)-2 {
				sb.WriteString(fmt.Sprintf("%s/", callers[i]))
				continue
			}

			if i == 0 && v == "vendor" {
				sb.WriteString("vendor/")
				continue
			}

			chars := []rune(v)
			if len(chars) > 0 {
				sb.WriteRune(chars[0])
				sb.WriteString("/")
			}
		}

		return sb.String()
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

// NewFromConfig returns logger based on config file
//
// log:
//   level: info
//   json-enabled: false
//   file:
//     enabled: true
//     path: ./logs/promo-engine.log
//     maxsize-mb: 10
//     maxage-days: 7
//     maxbackup-files: 2
//   batch:
//     enabled: false
//     max-lines: 1000
//     interval: 15ms
//
// then we can call using :
// v := viper.New()
// ... set v file configs, etc
//
// logger := log.NewFromConfig(v, "log")
// ..continue using logger.
func NewFromConfig(cfg kitconfig.KVStore, path string) (l *Logger, err error) {
	appName := cfg.GetString("name")

	logJSONFormat := cfg.GetBool(fmt.Sprintf("%s.json-enabled", path))
	logLevel := cfg.GetString(fmt.Sprintf("%s.level", path))
	logFilePath := cfg.GetString(fmt.Sprintf("%s.file.path", path))
	logFileMaxSize := cfg.GetInt(fmt.Sprintf("%s.file.maxsize-mb", path))
	logFileMaxAge := cfg.GetInt(fmt.Sprintf("%s.file.maxage-days", path))
	logFileMaxBackups := cfg.GetInt(fmt.Sprintf("%s.file.maxbackup-files", path))
	logFileEnabled := cfg.GetBool(fmt.Sprintf("%s.file.enabled", path))

	logBatchCfg := &BatchConfig{
		MaxLines: cfg.GetInt(fmt.Sprintf("%s.batch.max-lines", path)),
		Enabled:  cfg.GetBool(fmt.Sprintf("%s.batch.enabled", path)),
		Interval: cfg.GetDuration(fmt.Sprintf("%s.batch.interval", path)),
	}

	var fileLogger *lumberjack.Logger

	if logFileEnabled {
		logFile, err := os.OpenFile(filepath.Clean(logFilePath), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to open logfile %s", logFilePath)
		}

		fileLogger = &lumberjack.Logger{
			Filename:   logFile.Name(),
			MaxSize:    logFileMaxSize,
			MaxAge:     logFileMaxAge,
			MaxBackups: logFileMaxBackups,
			LocalTime:  true,
			Compress:   true,
		}
	}

	if logJSONFormat {
		// Use JSONFormatter logger for non development environment
		l = NewLogger(GetLevelFromString(logLevel), appName, fileLogger, logBatchCfg)
	} else {
		l = NewDevLogger(fileLogger, logBatchCfg)
	}

	return l, nil
}
