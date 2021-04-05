package springcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// ErrStatusCode ...
var ErrStatusCode = errors.New("springcloud: invalid status code response")

// Client interacts with remote config.
type Client struct {
	netClient *http.Client
	url       string
}

// AppConfig is app identifier in springcloud remote config.
type AppConfig struct {
	Name    string
	Profile string
	Branch  string
}

// NewRemoteConfigClient returns new springcloud config client.
func NewRemoteConfigClient(c *http.Client, url string) *Client {
	return &Client{
		netClient: c,
		url:       url,
	}
}

// structs having same structure as response from spring cloud config.
type springcloudconfig struct {
	Name            string           `json:"name"`
	Profiles        []string         `json:"profiles"`
	Label           string           `json:"label"`
	Version         string           `json:"version"`
	Propertysources []propertysource `json:"propertysources"`
}

type propertysource struct {
	Name   string                 `json:"name"`
	Source map[string]interface{} `json:"source"`
}

// LoadViperConfig parse spring cloud config values to *viper.Viper instance
// config source will be taken from <url>/<app-name>/<profile>/<branch>.
func (c *Client) LoadViperConfig(ctx context.Context, viper *viper.Viper, appCfg AppConfig) error {
	url := fmt.Sprintf("%s/%s/%s/%s", c.url, appCfg.Name, appCfg.Profile, appCfg.Branch)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return errors.Wrap(err, "gloudconfig: building request failed")
	}

	resp, err := c.netClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "gcloudconfig: http request failed")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Wrapf(ErrStatusCode, "gcloudconfig: status code %d", resp.StatusCode)
	}

	var cfg springcloudconfig

	if err = json.NewDecoder(resp.Body).Decode(&cfg); err != nil {
		return errors.Wrap(err, "gcloudconfig: parsing json response failed")
	}

	for key, value := range cfg.Propertysources[0].Source {
		viper.Set(key, value)
	}

	return nil
}
