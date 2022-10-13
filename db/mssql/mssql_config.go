package mssql

import (
	"database/sql"
	"fmt"

	"github.com/adipurnama/go-toolkit/config"
	"github.com/adipurnama/go-toolkit/db"
)

// NewFromConfig ...
func NewFromConfig(cfg config.KVStore, host, path string) (*sql.DB, error) {
	connOpt := db.DefaultConnectionOption()

	if maxIdle := cfg.GetInt(fmt.Sprintf("%s.conn.max-idle", path)); maxIdle > 0 {
		connOpt.MaxIdle = maxIdle
	}

	if maxLifetime := cfg.GetDuration(fmt.Sprintf("%s.conn.max-lifetime", path)); maxLifetime > 0 {
		connOpt.MaxLifetime = maxLifetime
	}

	if connTimeout := cfg.GetDuration(fmt.Sprintf("%s.conn.timeout", path)); connTimeout > 0 {
		connOpt.ConnectTimeout = connTimeout
	}

	if keepAlive := cfg.GetDuration(fmt.Sprintf("%s.conn.keep-alive-interval", path)); keepAlive > 0 {
		connOpt.KeepAliveCheckInterval = keepAlive
	}

	opt, err := db.NewDatabaseOption(
		host,
		cfg.GetInt(fmt.Sprintf("%s.port", path)),
		cfg.GetString(fmt.Sprintf("%s.username", path)),
		cfg.GetString(fmt.Sprintf("%s.password", path)),
		cfg.GetString(fmt.Sprintf("%s.schema", path)),
		connOpt,
	)
	if err != nil {
		return nil, err
	}

	return NewMsSQLDatabase(opt)
}
