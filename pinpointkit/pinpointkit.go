package pinpointkit

import (
	"fmt"

	"github.com/pinpoint-apm/pinpoint-go-agent"
	"github.com/pkg/errors"

	"github.com/adipurnama/go-toolkit/config"
)

const maxLengthAppName = 24

// Options ...
type Options struct {
	AppName string
	Env     string
	Host    string
}

// APMPinpointConfigs sets options for tracing.
type APMPinpointConfigs func(*Options)

// NewAgent returns pinpoint-agent
// configs appName and env will display on pinpoint web
// will show env-appName on pinpoint.
func NewAgent(c ...APMPinpointConfigs) (pinpoint.Agent, error) {
	conf := Options{
		AppName: "app_name",
		Env:     "staging",
		Host:    "localhost",
	}

	for _, c := range c {
		c(&conf)
	}

	appName := fmt.Sprintf("%s-%s", conf.Env, conf.AppName)

	if len(appName) > maxLengthAppName {
		appName = appName[0:maxLengthAppName]
	}

	pOpts := []pinpoint.ConfigOption{
		pinpoint.WithAgentId(appName),
		pinpoint.WithAppName(appName),
		pinpoint.WithCollectorHost(conf.Host),
	}

	pCfg, err := pinpoint.NewConfig(pOpts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed create new pinpoint config")
	}

	agent, err := pinpoint.NewAgent(pCfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed create new pinpoint agent")
	}

	return agent, nil
}

// WithOptionsFromConfig add specific name, env and host from cloud-config
// e.g
//
// app:
//
//	pinpoint:
//	 app-name: "app-name"
//	 environment: "staging, production"
//	 host: "".
func WithOptionsFromConfig(cfg config.KVStore, path string) APMPinpointConfigs {
	return func(c *Options) {
		if cfg.IsSet(fmt.Sprintf("%s.app-name", path)) {
			c.AppName = cfg.GetString(fmt.Sprintf("%s.app-name", path))
		}

		if cfg.IsSet(fmt.Sprintf("%s.environment", path)) {
			c.Env = cfg.GetString(fmt.Sprintf("%s.environment", path))
		}

		if cfg.IsSet(fmt.Sprintf("%s.host", path)) {
			c.Host = cfg.GetString(fmt.Sprintf("%s.host", path))
		}
	}
}

// WithOptions add specific name, env and host from struct internal options.
func WithOptions(opt Options) APMPinpointConfigs {
	return func(o *Options) {
		o.AppName = opt.AppName
		o.Env = opt.Env
		o.Host = opt.Host
	}
}
