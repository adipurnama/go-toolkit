package springcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var (
	// ErrStatusCode ...
	ErrStatusCode = errors.New("springcloud: invalid status code response")
	// ErrConfigNotFound ...
	ErrConfigNotFound = errors.New("springcloud: config not found at config server")
)

type (
	// Client interacts with remote config.
	Client struct {
		netClient *http.Client
	}

	appConfig struct {
		ConfigPath string `envconfig:"SPRING_CLOUD_CONFIG_PATH" required:"true"`
		ConfigURL  string `envconfig:"SPRING_CLOUD_CONFIG_URL" required:"true"`
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

func (cfg appConfig) confingEndpoint() string {
	cfg.ConfigURL = strings.TrimSuffix(cfg.ConfigURL, "/")
	cfg.ConfigPath = strings.TrimSuffix(cfg.ConfigPath, "/")
	cfg.ConfigPath = strings.TrimPrefix(cfg.ConfigPath, "/")

	return fmt.Sprintf("%s/%s", cfg.ConfigURL, cfg.ConfigPath)
}

// NewRemoteConfigClient returns new springcloud config client.
func NewRemoteConfigClient(c *http.Client) *Client {
	return &Client{
		netClient: c,
	}
}

// LoadViperConfig parse spring cloud config values to *viper.Viper instance
// config source will be taken from <url>/<app-name>/<profile>/<branch>.
func (c *Client) LoadViperConfig(ctx context.Context, viper *viper.Viper) error {
	var appCfg appConfig

	err := envconfig.Process("", &appCfg)
	if err != nil {
		return errors.Wrap(err, "springcloud: error parsing cloud config")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, appCfg.confingEndpoint(), nil)
	if err != nil {
		return errors.Wrap(err, "gloudconfig: building request failed")
	}

	resp, err := c.netClient.Do(req)
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
		return errors.Wrapf(ErrConfigNotFound, "config url %s", appCfg.confingEndpoint())
	}

	for key, value := range cfg.Propertysources[0].Source {
		viper.Set(key, value)
	}

	return nil
}
