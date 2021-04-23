package echokit

import (
	"fmt"

	"github.com/spf13/viper"
)

// NewViperRuntimeConfig returns *RuntimeConfig based on viper configuration
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
// call using `echokit.NewViperRuntimeConfig(v, "restapi")`.
func NewViperRuntimeConfig(v *viper.Viper, path string) *RuntimeConfig {
	r := RuntimeConfig{}

	r.Port = v.GetInt(fmt.Sprintf("%s.port", path))
	r.RequestTimeoutConfig = &TimeoutConfig{
		Timeout: v.GetDuration(fmt.Sprintf("%s.request-timeout", path)),
	}
	r.ShutdownTimeoutDuration = v.GetDuration(fmt.Sprintf("%s.shutdown.timeout-duration", path))
	r.ShutdownWaitDuration = v.GetDuration(fmt.Sprintf("%s.shutdown.wait-duration", path))
	r.HealthCheckPath = v.GetString(fmt.Sprintf("%s.healthcheck-path", path))
	r.InfoCheckPath = v.GetString(fmt.Sprintf("%s.info-path", path))

	return &r
}
