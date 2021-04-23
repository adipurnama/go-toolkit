package main

import (
	// "github.com/adipurnama/go-toolkit/db/postgres".
	"github.com/adipurnama/go-toolkit/log"
	"github.com/adipurnama/go-toolkit/runtimekit"
	"github.com/adipurnama/go-toolkit/springcloud"
	"github.com/adipurnama/go-toolkit/web/httpclient"

	"github.com/spf13/viper"
)

// test run using 'make run-springconfig' command
// make sure localhost spring cloud config is already running.
//  `docker-compose up -d spring-cloud-config`
func main() {
	appCtx, cancel := runtimekit.NewRuntimeContext()
	defer cancel()

	httpclient := httpclient.NewStdHTTPClient()
	configClient := springcloud.NewRemoteConfigClient(httpclient)
	v := viper.New()

	err := configClient.LoadViperConfig(appCtx, v)
	if err != nil {
		log.FromCtx(appCtx).Error(err, "load remote config failed")

		return
	}

	// db, err := postgres.NewFromViperFileConfig(v, "postgres.primary")
	// if err != nil {
	// 	log.FromCtx(appCtx).Error(err, "load db from viper config failed")

	// 	return
	// }

	// print all remote config key-values
	log.FromCtx(appCtx).Info("load remote config success")

	// setup logging
	isServerMode := v.GetBool("log.json-enabled")

	log.Println("json enabled ", isServerMode)

	if isServerMode {
		// production mode - json
		_ = log.NewLogger(log.LevelDebug, "springcloud_cfg_app", nil, nil).Set()
	} else {
		// development mode - logfmt
		_ = log.NewDevLogger(nil, nil).Set()
	}

	for _, k := range v.AllKeys() {
		log.FromCtx(appCtx).Info("===>", k, v.Get(k))
	}
}
