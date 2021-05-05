package echokit

import (
	"fmt"

	"github.com/adipurnama/go-toolkit/config"
)

// NewRuntimeConfig returns *RuntimeConfig based on viper configuration
// with layout:
//
// restapi:
//   port: 8088
//   request-timeout: 10s
//   healthcheck-path: /health/info
//   info-path: /actuator/info
//   shutdown:
//     wait-duration: 3s
//     timeout-duration: 5s
//
// call using `echokit.NewRuntimeConfig(v, "restapi")`.
func NewRuntimeConfig(cfg config.KVStore, path string) *RuntimeConfig {
	r := RuntimeConfig{}

	r.Port = cfg.GetInt(fmt.Sprintf("%s.port", path))
	r.RequestTimeoutConfig = &TimeoutConfig{
		Timeout: cfg.GetDuration(fmt.Sprintf("%s.request-timeout", path)),
	}
	r.ShutdownTimeoutDuration = cfg.GetDuration(fmt.Sprintf("%s.shutdown.timeout-duration", path))
	r.ShutdownWaitDuration = cfg.GetDuration(fmt.Sprintf("%s.shutdown.wait-duration", path))
	r.HealthCheckPath = cfg.GetString(fmt.Sprintf("%s.healthcheck-path", path))
	r.InfoCheckPath = cfg.GetString(fmt.Sprintf("%s.info-path", path))

	return &r
}
