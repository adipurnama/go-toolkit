package main

import (
	// "github.com/adipurnama/go-toolkit/db/postgres".

	"fmt"
	"net/http"

	"github.com/adipurnama/go-toolkit/log"
	"github.com/adipurnama/go-toolkit/runtimekit"
	"github.com/adipurnama/go-toolkit/springcloud"
	"github.com/adipurnama/go-toolkit/web/httpclient"
	"github.com/labstack/echo/v4"
)

// test run using 'make run-springconfig' command
// make sure localhost spring cloud config is already running.
//  `docker compose up -d spring-cloud-config`
func main() {
	appCtx, cancel := runtimekit.NewRuntimeContext()
	defer cancel()

	httpclient := httpclient.NewStdHTTPClient()
	config := springcloud.NewRemoteConfig(httpclient)

	err := config.Load(appCtx)
	if err != nil {
		log.FromCtx(appCtx).Error(err, "load remote config failed")

		return
	}

	// print all remote config key-values
	log.FromCtx(appCtx).Info("load remote config success")

	// setup logging
	isServerMode := config.GetBool("log.json-enabled")

	log.Println("json enabled ", isServerMode)

	if isServerMode {
		// production mode - json
		_ = log.NewLogger(log.LevelDebug, "springcloud_cfg_app", nil, nil).Set()
	} else {
		// development mode - logfmt
		_ = log.NewDevLogger(nil, nil).Set()
	}

	for _, k := range config.AllKeys() {
		log.FromCtx(appCtx).Info("===>", k, config.Get(k))
	}

	// create db, pubsub, etc from config properties

	// db, err := postgres.NewFromViperFileConfig(v, "postgres.primary")
	// if err != nil {
	// 	log.FromCtx(appCtx).Error(err, "load db from viper config failed")

	// 	return
	// }

	e := echo.New()
	e.GET("/values", getValues(config))

	go func() {
		if err := e.Start(":8081"); err != http.ErrServerClosed {
			log.FromCtx(appCtx).Error(err, "failed serving http server")
		}
	}()

	<-appCtx.Done()

	log.FromCtx(appCtx).Info("Bye.")
}

func getValues(cfg *springcloud.RemoteConfig) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		key := ctx.QueryParam("key")

		if !cfg.IsSet(key) {
			return ctx.String(http.StatusBadRequest, fmt.Sprintf("key `%s` doesn't exists", key))
		}

		return ctx.String(http.StatusOK, fmt.Sprintf("%v", cfg.Get(key)))
	}
}
