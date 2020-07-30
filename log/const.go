package log

// Level is log output level
type Level int

const (
	// LevelDisabled level
	LevelDisabled Level = -1
	// LevelDebug level
	LevelDebug Level = 0
	// LevelInfo level
	LevelInfo Level = 1
	// LevelWarn level
	LevelWarn Level = 2
	// LevelError level
	LevelError Level = 3

	cfgSkipCallerCount              = 4
	cfgDefaultStdLogSkipCallerCount = 2
)

var (
	loggerCtxKey contextKey = contextKey{name: "internal-ctx-log"}
)
