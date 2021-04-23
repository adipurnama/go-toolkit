package repository

import (
	"context"
	"errors"
	"time"
)

var (
	// ErrNoRows no result found.
	ErrNoRows = errors.New("repository fakesql: no rows found")
	// ErrConnectTimeout simulates timeout to db call.
	ErrConnectTimeout = errors.New("repository fakesql: context deadline exceeded")
)

// DummyDB as a buggy *sql.DB.
type DummyDB struct {
}

// QueryRowsCtx sql.DB impl.
func (db *DummyDB) QueryRowsCtx(_ context.Context, _ string) error {
	ok := getRandomBool()
	if !ok {
		return ErrNoRows
	}

	return nil
}

// QueryRowCtx sql.DB impl.
func (db *DummyDB) QueryRowCtx(_ context.Context, _ string) (int, error) {
	ok := getRandomBool()
	if !ok {
		return 0, ErrConnectTimeout
	}

	dummyID := 10

	return dummyID, nil
}

// Ping sql.DB impl.
func (db *DummyDB) Ping(_ context.Context) error {
	ok := getRandomBool()
	if !ok {
		return ErrConnectTimeout
	}

	return nil
}

func getRandomBool() bool {
	now := time.Now()
	return now.Second()%2 == 0
}
