package rediskit

import (
	"fmt"

	"github.com/adipurnama/go-toolkit/db"
	goredis "github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

// NewFromViperFileConfig returns redis *redis.Client instance from yaml config file
//
// redis:
//   primary:
//     username: <username>
//     password: "<password>"
//     host: mredis.aws.com
//     port: 6379
//     schema: 0
//     conn:
//       max-idle: 20
//       max-lifetime: 10m
//       timeout: 5m
//       max-open: 100
//
// then we can call using :
// v := viper.New()
// ... set v file configs, etc
//
// db, err := log.NewFromViperFileConfig(v, "redis.primary")
// ...continue using db.
func NewFromViperFileConfig(v *viper.Viper, path string) (*goredis.Client, error) {
	connOpt := db.DefaultConnectionOption()

	if maxIdle := v.GetInt(fmt.Sprintf("%s.conn.max-idle", path)); maxIdle > 0 {
		connOpt.MaxIdle = v.GetInt(fmt.Sprintf("%s.conn.max-idle", path))
	}

	if maxOpen := v.GetInt(fmt.Sprintf("%s.conn.max-open", path)); maxOpen > 0 {
		connOpt.MaxOpen = maxOpen
	}

	if maxLifetime := v.GetDuration(fmt.Sprintf("%s.conn.max-lifetime", path)); maxLifetime > 0 {
		connOpt.MaxLifetime = maxLifetime
	}

	if connTimeout := v.GetDuration(fmt.Sprintf("%s.conn.timeout", path)); connTimeout > 0 {
		connOpt.ConnectTimeout = connTimeout
	}

	opt, err := db.NewDatabaseOption(
		v.GetString(fmt.Sprintf("%s.host", path)),
		v.GetInt(fmt.Sprintf("%s.port", path)),
		v.GetString(fmt.Sprintf("%s.username", path)),
		v.GetString(fmt.Sprintf("%s.password", path)),
		v.GetString(fmt.Sprintf("%s.schema", path)),
		connOpt,
	)
	if err != nil {
		return nil, err
	}

	return NewRedisConnection(opt)
}
