// Package springcloud helps interact with spring cloud remote config
package springcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/adipurnama/go-toolkit/log"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var (
	// ErrStatusCode returned when call to springcloud remote config returns HTTP other than 2xx.
	ErrStatusCode = errors.New("springcloud: invalid status code response")
	// ErrConfigNotFound returned when profiles doesn't exists in remote config.
	ErrConfigNotFound = errors.New("springcloud: config not found at config server")
)

const (
	// secret-manager value should be in format `{sm}my-secret-key`
	// will look for `my-secret-key` value inside gcp secretmanager.
	gcpSecretManagerOffset = 4

	cfgKeySpringConfigRefreshInterval = "springcloud-config.refresh-interval"
)

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
	cfg.ConfigURL = strings.TrimSuffix(cfg.ConfigURL, "/")

	// check for local file case
	if strings.HasPrefix(cfg.ConfigURL, "file://") {
		return []string{cfg.ConfigURL}
	}

	// setup remote springconfig urls
	urls := []string{}

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
		log.FromCtx(ctx).Info("springcloud_config: auto-refresh is inactive")
		return nil
	}

	refreshInterval := c.cfg.GetDuration(cfgKeySpringConfigRefreshInterval)
	if refreshInterval == 0 {
		log.FromCtx(ctx).Info("springcloud_config: auto-refresh is inactive for zero interval", "interval", refreshInterval)
		return nil
	}

	log.FromCtx(ctx).Info("springcloud_config: run with auto-refresh config", "interval", refreshInterval)

	go func() {
		for {
			select {
			case <-time.After(refreshInterval):
				start := time.Now()

				errRefresh := c.loadRemoteConfigForViper(context.Background())
				if errRefresh != nil {
					log.FromCtx(ctx).Error(errRefresh, "refresh remote config failed")
					continue
				}

				log.FromCtx(ctx).Info("springcloud_config: config reloaded success",
					"elapsed_time_ms", time.Since(start).Milliseconds(),
				)
			case <-ctx.Done():
				log.FromCtx(ctx).
					Info("springcloud_config: context.Done: stopping springcloud config auto-refresh",
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
		return errors.Wrap(err, "springcloud_config: error parsing cloud config")
	}

	for _, urlStr := range appCfg.configEndpoints() {
		parsedURL, err := url.Parse(urlStr)
		if err != nil {
			return errors.Wrap(err, "springcloud_config: failed to parse url")
		}

		if parsedURL.Scheme == "file" {
			return c.applyViperFromLocalFile(ctx, parsedURL)
		}

		err = c.applyViperFromSpringRemoteURL(ctx, urlStr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *RemoteConfig) applyViperFromLocalFile(
	ctx context.Context,
	fileURL *url.URL,
) error {
	v := viper.New()
	path := strings.Replace(fileURL.String(), "file://", "", 1)

	path = strings.TrimSuffix(path, "/")
	log.Println("path", path)

	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return errors.Wrap(err, "springcloud_config: failed to read local file")
	}

	log.FromCtx(ctx).Info("meh", "keyvals", v.AllSettings())

	return c.applyKeyValues(ctx, v.AllSettings())
}

func (c *RemoteConfig) applyViperFromSpringRemoteURL(ctx context.Context, url string) (err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return errors.Wrap(err, "springcloud_config: building request failed")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "springcloud_config: http request failed")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Wrapf(ErrStatusCode, "springcloud_config: status code %d", resp.StatusCode)
	}

	var cfg springcloudconfig

	if err = json.NewDecoder(resp.Body).Decode(&cfg); err != nil {
		return errors.Wrap(err, "springcloud_config: parsing json response failed")
	}

	if len(cfg.Propertysources) == 0 {
		return errors.Wrapf(ErrConfigNotFound, "config_url=%s", url)
	}

	return c.applyKeyValues(ctx, cfg.Propertysources[0].Source)
}

func (c *RemoteConfig) applyKeyValues(ctx context.Context, keyValues map[string]interface{}) (err error) {
	var smClient *secretmanager.Client

	c.mu.Lock()

	defer func() {
		c.mu.Unlock()

		if smClient != nil {
			smClient.Close()
		}
	}()

	for key, value := range keyValues {
		strValue := fmt.Sprintf("%v", value)

		// set for plain key-values
		if !strings.HasPrefix(strValue, "{sm}") {
			c.cfg.Set(key, value)
			continue
		}

		// set for gcp secret manager values

		// create the secretmanager client if it hasn't initialized yet
		if smClient == nil {
			smClient, err = secretmanager.NewClient(ctx)
			if err != nil {
				return errors.Wrapf(
					err,
					"springcloud_config: failed to create secretmanager client for key %s", key,
				)
			}
		}

		// retrieve and apply secret manager value
		smValue, err := getSecretManagerValue(ctx, smClient, strValue[gcpSecretManagerOffset:])
		if err != nil {
			return err
		}

		c.cfg.Set(key, smValue)
	}

	return nil
}
