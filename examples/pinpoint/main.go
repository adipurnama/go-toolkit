package main

import (
	"github.com/adipurnama/go-toolkit/log"
	"github.com/adipurnama/go-toolkit/pinpointkit"
	"github.com/adipurnama/go-toolkit/runtimekit"
	"github.com/adipurnama/go-toolkit/springcloud"
	"github.com/adipurnama/go-toolkit/web/httpclient"
)

func main() {
	var err error

	ctx, cancel := runtimekit.NewRuntimeContext()

	defer func() {
		cancel()

		if err != nil {
			log.FromCtx(ctx).Error(err, "found error")
		}
	}()

	httpClient := httpclient.NewStdHTTPClient()
	springConfig := springcloud.NewRemoteConfig(httpClient)

	conf := pinpointkit.WithOptionsFromConfig(springConfig, "app.pinpoint")

	_, err = pinpointkit.NewAgent(conf)
	if err != nil {
		log.FromCtx(ctx).Error(err, "cannot init pinpoint agent")
	}
}
