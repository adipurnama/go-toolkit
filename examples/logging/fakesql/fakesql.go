package fakesql

import (
	"context"
	"errors"
)

// ErrNoRows no result found.
var ErrNoRows = errors.New("fakesql: no rows found")

// DB as a buggy *sql.DB.
type DB struct {
}

// QueryCtx sql.DB impl.
func (db *DB) QueryCtx(_ context.Context, _ string) error {
	return ErrNoRows
}
