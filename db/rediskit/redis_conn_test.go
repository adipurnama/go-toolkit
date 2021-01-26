package rediskit_test

import (
	"log"
	"strconv"
	"testing"

	"github.com/adipurnama/go-toolkit/db"
	"github.com/adipurnama/go-toolkit/db/rediskit"
	"github.com/alicebob/miniredis"
)

func TestNewRedisConnection(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	port, _ := strconv.Atoi(mr.Port())

	dbOpt, _ := db.NewDatabaseOption(mr.Host(), port, "", "", "", nil)

	client, err := rediskit.NewRedisConnection(dbOpt)
	if err != nil {
		t.Fatal("should return no error")
	}

	if client == nil {
		t.Fatal("should return valid redis client")
	}
}
