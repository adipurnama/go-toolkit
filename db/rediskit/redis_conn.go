package rediskit

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/adipurnama/go-toolkit/db"
	goredis "github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	apmgoredis "go.elastic.co/apm/module/apmgoredisv8"
)

// NewRedisConnection returns new redis client
// based on db Options.
func NewRedisConnection(option *db.Option) (*goredis.Client, error) {
	opts := goredis.Options{
		Addr:         fmt.Sprintf("%s:%d", option.Host, option.Port),
		Password:     option.Password,
		DialTimeout:  option.ConnectTimeout,
		MinIdleConns: option.ConnectionOption.MaxIdle,
		PoolSize:     option.ConnectionOption.MaxOpen,
		MaxConnAge:   option.ConnectionOption.MaxLifetime,
		PoolTimeout:  option.MaxLifetime,
	}

	dbID, err := strconv.Atoi(option.DatabaseName)
	if err == nil {
		opts.DB = dbID
	}

	rClient := goredis.NewClient(&opts)
	rClient.AddHook(apmgoredis.NewHook())

	_, err = rClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, errors.Wrap(err, "rediskit: failed to initiate redis PING")
	}

	log.Println("successfully connected to redis", opts.Addr)

	return rClient, nil
}
