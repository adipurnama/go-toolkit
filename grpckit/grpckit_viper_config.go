package grpckit

import (
	"fmt"

	"github.com/spf13/viper"
)

// NewViperRuntimeConfig returns *RuntimeConfig based on viper configuration
// with layout:
//
// grpc:
//   port: 8088
//   request-timeout: 10s
//   shutdown-wait-duration: 3s
//   reflection-enabled: true
//
// call using `grpckit.NewViperRuntimeConfig(v, "grpc")`.
func NewViperRuntimeConfig(v *viper.Viper, path string) *RuntimeConfig {
	r := RuntimeConfig{}

	r.Port = v.GetInt(fmt.Sprintf("%s.port", path))
	r.RequestTimeout = v.GetDuration(fmt.Sprintf("%s.request-timeout", path))
	r.ShutdownWaitDuration = v.GetDuration(fmt.Sprintf("%s.shutdown-wait-duration", path))
	r.EnableReflection = v.GetBool(fmt.Sprintf("%s.reflection-enabled", path))

	return &r
}
