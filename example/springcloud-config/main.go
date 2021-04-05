package main

import (
	"os"

	"github.com/adipurnama/go-toolkit/log"
	"github.com/adipurnama/go-toolkit/runtimekit"
	"github.com/adipurnama/go-toolkit/springcloud"
	"github.com/adipurnama/go-toolkit/web/httpclient"
	"github.com/spf13/viper"
)

func main() {
	remoteConfigURL := os.Getenv("SPRING_CLOUD_CONFIG_URL")
	if remoteConfigURL == "" {
		log.Fatal("env SPRING_CLOUD_CONFIG_URL cannot be empty")
	}

	appCtx, cancel := runtimekit.NewRuntimeContext()
	defer cancel()

	httpclient := httpclient.NewStdHTTPClient()
	configClient := springcloud.NewRemoteConfigClient(httpclient, remoteConfigURL)
	v := viper.New()

	// if we're using git-backed springcloud config,
	// we should have file 'config-debugger-development.yaml' at branch 'master'
	appCfg := springcloud.AppConfig{
		Name:    "config-debugger",
		Profile: "development",
		Branch:  "master",
	}

	err := configClient.LoadViperConfig(appCtx, v, appCfg)
	if err != nil {
		log.FromCtx(appCtx).Error(err, "load remote config failed")

		return
	}

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
