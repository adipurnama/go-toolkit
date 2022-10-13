package rediskit

import (
	"fmt"

	goredis "github.com/go-redis/redis/v8"

	"github.com/adipurnama/go-toolkit/config"
	"github.com/adipurnama/go-toolkit/db"
)

/*
NewFromConfig returns redis *redis.Client instance from yaml config file

	given config file contents:

		redis:
		  primary:
			username: <username>
			password: "<password>"
			host: mredis.aws.com
			port: 6379
			schema: 0
			conn:
			  max-idle: 20
			  max-lifetime: 10m
			  timeout: 5m
			  max-open: 100

	then we can call using :

		v := viper.New()
		... set v file configs, etc

		db, err := log.NewFromConfig(v, "redis.primary")
		...continue using db.
*/
func NewFromConfig(cfg config.KVStore, path string) (*goredis.Client, error) {
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

	return NewRedisConnection(opt)
}
