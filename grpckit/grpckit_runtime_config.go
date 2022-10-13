package grpckit

import (
	"fmt"

	"github.com/adipurnama/go-toolkit/config"
)

/*
NewRuntimeConfig returns *RuntimeConfig based on viper configuration
with layout:

	given config file contents:

		grpc:
		  port: 8088
		  request-timeout: 10s
		  shutdown-wait-duration: 3s
		  reflection-enabled: true

	call using `grpckit.NewRuntimeConfig(v, "grpc")`.
*/
func NewRuntimeConfig(cfg config.KVStore, path string) *RuntimeConfig {
	r := RuntimeConfig{}

	r.Port = cfg.GetInt(fmt.Sprintf("%s.port", path))
	r.RequestTimeout = cfg.GetDuration(fmt.Sprintf("%s.request-timeout", path))
	r.ShutdownWaitDuration = cfg.GetDuration(fmt.Sprintf("%s.shutdown-wait-duration", path))
	r.EnableReflection = cfg.GetBool(fmt.Sprintf("%s.reflection-enabled", path))

	return &r
}
