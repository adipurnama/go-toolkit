package log

const (
	// Disabled level
	Disabled = -1
	// LevelDebug level
	LevelDebug = iota
	// LevelInfo level
	LevelInfo
	// LevelWarn level
	LevelWarn
	// LevelError level
	LevelError
)

var (
	loggerCtxKey contextKey = contextKey{name: "internal-ctx-log"}
)
