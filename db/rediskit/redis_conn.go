package rediskit

import (
	"context"
	"fmt"
	"log"

	"github.com/adipurnama/go-toolkit/db"
	goredis "github.com/go-redis/redis/v8"
)

// NewRedisConnection returns new redis client
// based on db Options
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

	rClient := goredis.NewClient(&opts)

	_, err := rClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	log.Println("successfully connected to redis")

	return rClient, nil
}
