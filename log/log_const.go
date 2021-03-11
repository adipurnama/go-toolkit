package log

import "strings"

// Level is log output level.
type Level int

const (
	// LevelDisabled level.
	LevelDisabled Level = -1
	// LevelDebug level.
	LevelDebug Level = 0
	// LevelInfo level.
	LevelInfo Level = 1
	// LevelWarn level.
	LevelWarn Level = 2
	// LevelError level.
	LevelError Level = 3

	cfgSkipCallerCount              = 4
	cfgDefaultStdLogSkipCallerCount = 2
)

var (
	loggerCtxKey contextKey = contextKey{name: "internal-ctx-log"}
)

func levelString(l Level) string {
	switch l {
	case LevelError:
		return "ERROR"
	case LevelWarn:
		return "WARN"
	case LevelInfo:
		return "INFO"
	default:
		return "DEBUG"
	}
}

// GetLevelFromString return error level based on config string.
func GetLevelFromString(level string) Level {
	switch strings.ToLower(level) {
	case "warn":
		return LevelWarn
	case "debug":
		return LevelDebug
	case "error":
		return LevelError
	default:
		return LevelInfo
	}
}
