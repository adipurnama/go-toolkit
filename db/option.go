package db

import "errors"

// Option - database option.
type Option struct {
	Host         string
	Port         int
	Username     string
	Password     string
	DatabaseName string
}

var errInvalidDbSource = errors.New("invalid datasource host | port")

// NewDatabaseOption - default factory method.
func NewDatabaseOption(host string, port int, username, password, dbName string) (*Option, error) {
	if host == "" || port == 0 {
		return nil, errInvalidDbSource
	}

	return &Option{
		Host:         host,
		Port:         port,
		Username:     username,
		Password:     password,
		DatabaseName: dbName,
	}, nil
}
