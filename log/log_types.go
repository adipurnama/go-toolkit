package log

import (
	"time"

	"github.com/rs/zerolog"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var cfg config

// Logger is structured leveled logger.
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
	logFmt     bool
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

	batchCfg *BatchConfig

	isDevelopment bool
}

// ErrFunc is any function which takes no argument and possibly returns error
// e.g. tx.Rollback(), resp.Body.Close().
type ErrFunc func() error

type contextKey struct {
	name string
}

// BatchConfig  is configuration for async log batch writer.
type BatchConfig struct {
	Enabled  bool
	Interval time.Duration
	MaxLines int
}

// String returns formatted context key identifier.
func (k *contextKey) String() string {
	return "mw-" + k.name
}

// setup name and static fields.
// Each new instance of logger will always append these
// key-value pairs to the output and name if it is not empty.
// These values cannot be modified after they are configured.
func setup(level Level, name string, fileLogger *lumberjack.Logger, isDevelopment bool, batchCfg *BatchConfig, stfields []interface{}) {
	if cfg.configured {
		return
	}

	cfg.isDevelopment = isDevelopment
	cfg.name = name
	cfg.stfields = append(cfg.stfields, stfields...)
	cfg.configured = true
	cfg.fileLogger = fileLogger
	cfg.batchCfg = batchCfg
	cfg.level = level
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
