package log

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

// NewFromViperFileConfig returns logger based on config file
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
// logger := log.NewFromViperFileConfig(v, "log")
// ..continue using logger.
func NewFromViperFileConfig(v *viper.Viper, path string) (l *Logger, err error) {
	appName := v.GetString("name")

	logJSONFormat := v.GetBool(fmt.Sprintf("%s.json-enabled", path))
	logLevel := v.GetString(fmt.Sprintf("%s.level", path))
	logFilePath := v.GetString(fmt.Sprintf("%s.file.path", path))
	logFileMaxSize := v.GetInt(fmt.Sprintf("%s.file.maxsize-mb", path))
	logFileMaxAge := v.GetInt(fmt.Sprintf("%s.file.maxage-days", path))
	logFileMaxBackups := v.GetInt(fmt.Sprintf("%s.file.maxbackup-files", path))
	logFileEnabled := v.GetBool(fmt.Sprintf("%s.file.enabled", path))

	logBatchCfg := &BatchConfig{
		MaxLines: v.GetInt(fmt.Sprintf("%s.batch.max-lines", path)),
		Enabled:  v.GetBool(fmt.Sprintf("%s.batch.enabled", path)),
		Interval: v.GetDuration(fmt.Sprintf("%s.batch.interval", path)),
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
