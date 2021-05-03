// Package springcloud helps interact with spring cloud remote config
package springcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/adipurnama/go-toolkit/log"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

var (
	// ErrStatusCode returned when call to springcloud remote config returns HTTP other than 2xx.
	ErrStatusCode = errors.New("springcloud: invalid status code response")
	// ErrConfigNotFound returned when profiles doesn't exists in remote config.
	ErrConfigNotFound = errors.New("springcloud: config not found at config server")
)

const cfgKeySpringConfigRefreshInterval = "springcloud-config.refresh-interval"

type (
	appConfig struct {
		ConfigPaths []string `envconfig:"SPRING_CLOUD_CONFIG_PATHS" required:"true"`
		ConfigURL   string   `envconfig:"SPRING_CLOUD_CONFIG_URL" required:"true"`
	}

	// structs having same structure as response from spring cloud config.
	springcloudconfig struct {
		Name            string           `json:"name"`
		Profiles        []string         `json:"profiles"`
		Label           string           `json:"label"`
		Version         string           `json:"version"`
		Propertysources []propertysource `json:"propertysources"`
	}

	propertysource struct {
		Name   string                 `json:"name"`
		Source map[string]interface{} `json:"source"`
	}
)

func (cfg appConfig) configEndpoints() []string {
	urls := []string{}

	cfg.ConfigURL = strings.TrimSuffix(cfg.ConfigURL, "/")

	for _, v := range cfg.ConfigPaths {
		path := strings.TrimSuffix(v, "/")
		path = strings.TrimPrefix(path, "/")

		url := fmt.Sprintf("%s/%s", cfg.ConfigURL, path)

		urls = append(urls, url)
	}

	return urls
}

// Load key-values from spring cloud config profile(s)
// if value `springcloud-config.refresh-interval` exists, preiodically refresh the values
// based on interval set.
func (c *RemoteConfig) Load(ctx context.Context) error {
	err := c.loadRemoteConfigForViper(ctx)
	if err != nil {
		return err
	}

	if !c.cfg.IsSet(cfgKeySpringConfigRefreshInterval) {
		log.FromCtx(ctx).Info("viper springcloud-config auto-refresh is inactive")
		return nil
	}

	refreshInterval := c.cfg.GetDuration(cfgKeySpringConfigRefreshInterval)
	if refreshInterval == 0 {
		log.FromCtx(ctx).Info("viper springcloud-config auto-refresh is inactive for zero interval", "interval", refreshInterval)
		return nil
	}

	log.FromCtx(ctx).Info("run viper springcloud-config auto-refresh", "interval", refreshInterval)

	go func() {
		ticker := time.NewTicker(refreshInterval)

		for {
			select {
			case <-ticker.C:
				start := time.Now()

				errRefresh := c.loadRemoteConfigForViper(context.Background())
				if errRefresh != nil {
					log.FromCtx(ctx).Error(errRefresh, "refresh remote config failed")
					continue
				}

				log.FromCtx(ctx).Info("config re-loaded",
					"elapsed_time_ms", time.Since(start).Milliseconds(),
				)
			case <-ctx.Done():
				log.FromCtx(ctx).
					Info("context.Done: stopping springcloud config auto-refresh",
						"error", ctx.Err(),
					)

				return
			}
		}
	}()

	return nil
}

func (c *RemoteConfig) loadRemoteConfigForViper(ctx context.Context) error {
	var appCfg appConfig

	err := envconfig.Process("", &appCfg)
	if err != nil {
		return errors.Wrap(err, "springcloud: error parsing cloud config")
	}

	for _, url := range appCfg.configEndpoints() {
		err = c.applyViperFromSpringRemoteURL(ctx, url)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *RemoteConfig) applyViperFromSpringRemoteURL(ctx context.Context, url string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return errors.Wrap(err, "springcloud: building request failed")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "springcloud: http request failed")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Wrapf(ErrStatusCode, "springcloud: status code %d", resp.StatusCode)
	}

	var cfg springcloudconfig

	if err = json.NewDecoder(resp.Body).Decode(&cfg); err != nil {
		return errors.Wrap(err, "springcloud: parsing json response failed")
	}

	if len(cfg.Propertysources) == 0 {
		return errors.Wrapf(ErrConfigNotFound, "config_url=%s", url)
	}

	c.mu.Lock()

	for key, value := range cfg.Propertysources[0].Source {
		c.cfg.Set(key, value)
	}

	c.mu.Unlock()

	return nil
}
