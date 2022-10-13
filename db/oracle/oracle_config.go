package oracle

import (
	"database/sql"
	"fmt"

	"github.com/adipurnama/go-toolkit/config"
	"github.com/adipurnama/go-toolkit/db"
)

/*
NewFromConfig returns postgres-based *sqlx.DB instance from yaml config file

	given config file contents:

		db:
		  primary:
			username: <username>
			password: "<password>"
			host: mydb.ap-southeast-1.rds.amazonaws.com
			port: 5432
			schema: my_db_schema
			conn:
			  max-idle: 20
			  max-lifetime: 10m
			  timeout: 5m
			  max-open: 100
			  keep-alive-interval: 30s

	then we can call using :

		v := viper.New()
		... set v file configs, etc

		db, err := log.NewFromConfig(v, "db.primary")
		...continue using db.
*/
func NewFromConfig(cfg config.KVStore, path string) (*sql.DB, error) {
	connOpt := db.DefaultConnectionOption()

	if maxIdle := cfg.GetInt(fmt.Sprintf("%s.conn.max-idle", path)); maxIdle > 0 {
		connOpt.MaxIdle = cfg.GetInt(fmt.Sprintf("%s.conn.max-idle", path))
	}

	if maxOpen := cfg.GetInt(fmt.Sprintf("%s.conn.max-open", path)); maxOpen > 0 {
		connOpt.MaxOpen = maxOpen
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
		cfg.GetString(fmt.Sprintf("%s.host", path)),
		cfg.GetInt(fmt.Sprintf("%s.port", path)),
		cfg.GetString(fmt.Sprintf("%s.username", path)),
		cfg.GetString(fmt.Sprintf("%s.password", path)),
		cfg.GetString(fmt.Sprintf("%s.schema", path)),
		connOpt,
	)
	if err != nil {
		return nil, err
	}

	return NewOracleDatabase(opt)
}
