package log

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

type (
	// Level - log level
	Level int
)

const (
	// LevelPanic - log Panic
	LevelPanic Level = iota + 1

	// LevelFatal - log Fatal
	LevelFatal

	// LevelError - log Error
	LevelError

	// LevelDebug - log Debug
	LevelDebug

	// LevelWarn - log Warning
	LevelWarn

	// LevelInfo - log Info
	LevelInfo
)

/*
SetupWithLogfmtOutput will setup global logger with logfmt output (https://www.brandur.org/logfmt)
example: 2019-07-23 10:57:18  INFO   **request completed** method=POST path=/cached-member-data
*/
func SetupWithLogfmtOutput(loc *time.Location) {
	output := zerolog.ConsoleWriter{Out: os.Stdout}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf(" %-6s", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("**%s**", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s=", i)
	}
	output.FormatErrFieldName = func(i interface{}) string {
		return "error="
	}
	output.FormatErrFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%+v", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	// output.FormatCaller = func(i interface{}) string {
	// 	return fmt.Sprintf("caller=%v", i)
	// }
	output.FormatTimestamp = func(i interface{}) string {
		ts, ok := i.(string)
		if !ok {
			return "nok error"
		}

		t, err := time.Parse(time.RFC3339, ts)
		if err != nil {
			return "error"
		}

		t = t.In(loc)

		return fmt.Sprintf("%d-%.2d-%.2d %.2d:%.2d:%.2d",
			t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	}
	zlog.Logger = zlog.Output(output)
}

// SetGlobalLevel - set logging level.
func SetGlobalLevel(l Level) {
	switch l {
	case LevelPanic:
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case LevelFatal:
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case LevelError:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case LevelDebug:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case LevelWarn:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case LevelInfo:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
