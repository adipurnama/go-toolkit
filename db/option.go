// Package db for interacting with database
package db

import (
	"errors"
	"time"
)

const (
	defaultMaxOpen        = 100
	defaultMaxLifetime    = 10 * time.Minute
	defaultMaxIdle        = 5
	defaultConnectTimeout = 10 * time.Second
)

// Option - database option.
type Option struct {
	Host         string
	Port         int
	Username     string
	Password     string
	DatabaseName string
	*ConnectionOption
}

// ConnectionOption is db connection option
type ConnectionOption struct {
	MaxIdle        int
	MaxLifetime    time.Duration
	MaxOpen        int
	ConnectTimeout time.Duration
}

// DefaultConnectionOption returns sensible conn setting
func DefaultConnectionOption() *ConnectionOption {
	return &ConnectionOption{
		MaxIdle:        defaultMaxIdle,
		MaxOpen:        defaultMaxOpen,
		MaxLifetime:    defaultMaxLifetime,
		ConnectTimeout: defaultConnectTimeout,
	}
}

var errInvalidDBSource = errors.New("invalid datasource host | port")

// NewDatabaseOption - default factory method.
func NewDatabaseOption(host string, port int, username, password, dbName string, conn *ConnectionOption) (*Option, error) {
	if host == "" || port == 0 {
		return nil, errInvalidDBSource
	}

	if conn == nil {
		conn = DefaultConnectionOption()
	}

	return &Option{
		Host:             host,
		Port:             port,
		Username:         username,
		Password:         password,
		DatabaseName:     dbName,
		ConnectionOption: conn,
	}, nil
}
