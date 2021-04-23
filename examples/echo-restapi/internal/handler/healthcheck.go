package handler

import (
	"context"

	"github.com/adipurnama/go-toolkit/echokit"
	"github.com/adipurnama/go-toolkit/examples/echo-restapi/internal/repository"
)

// HealthCheck health check handler.
func HealthCheck(db *repository.DummyDB) echokit.HealthCheckFunc {
	return func(ctx context.Context) error {
		err := db.Ping(ctx)
		return err
	}
}
