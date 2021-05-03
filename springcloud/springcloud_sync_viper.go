package springcloud

import (
	"net/http"
	"sync"
	"time"

	"github.com/spf13/viper"
)

// RemoteConfig wraps thread-safe *viper.Viper key-values from springcloud remote-config.
type RemoteConfig struct {
	cfg    *viper.Viper
	mu     *sync.Mutex
	client *http.Client
}

// NewRemoteConfig create *RemoteConfig with given existing *http.Client.
func NewRemoteConfig(client *http.Client) *RemoteConfig {
	return &RemoteConfig{
		cfg:    viper.New(),
		client: client,
		mu:     &sync.Mutex{},
	}
}

// ============ thread-safe *viper.Viper wrapper ============

// IsSet ...
func (c *RemoteConfig) IsSet(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.cfg.IsSet(key)
}

// Get ...
func (c *RemoteConfig) Get(key string) interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.cfg.Get(key)
}

// AllKeys ...
func (c *RemoteConfig) AllKeys() []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.cfg.AllKeys()
}

// GetBool ...
func (c *RemoteConfig) GetBool(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.cfg.GetBool(key)
}

// ViperValue ...
func (c *RemoteConfig) ViperValue() viper.Viper {
	c.mu.Lock()
	defer c.mu.Unlock()

	return *c.cfg
}

// =========== Time ================

// GetTime ...
func (c *RemoteConfig) GetTime(key string) time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.cfg.GetTime(key)
}

// GetDuration ...
func (c *RemoteConfig) GetDuration(key string) time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.cfg.GetDuration(key)
}

// =========== Int ================

// GetInt ...
func (c *RemoteConfig) GetInt(key string) int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.cfg.GetInt(key)
}

// GetInt32 ...
func (c *RemoteConfig) GetInt32(key string) int32 {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.cfg.GetInt32(key)
}

// GetInt64 ...
func (c *RemoteConfig) GetInt64(key string) int64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.cfg.GetInt64(key)
}

// GetIntSlice ...
func (c *RemoteConfig) GetIntSlice(key string) []int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.cfg.GetIntSlice(key)
}

// =========== String ================

// GetString ...
func (c *RemoteConfig) GetString(key string) string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.cfg.GetString(key)
}

// GetStringSlice ...
func (c *RemoteConfig) GetStringSlice(key string) []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.cfg.GetStringSlice(key)
}

// GetStringMap ...
func (c *RemoteConfig) GetStringMap(key string) map[string]interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.cfg.GetStringMap(key)
}

// GetStringMapStringSlice ...
func (c *RemoteConfig) GetStringMapStringSlice(key string) map[string][]string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.cfg.GetStringMapStringSlice(key)
}

// ============== Float ================

// GetFloat64 ...
func (c *RemoteConfig) GetFloat64(key string) float64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.cfg.GetFloat64(key)
}

// ============== Uint ================

// GetUint ...
func (c *RemoteConfig) GetUint(key string) uint {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.cfg.GetUint(key)
}

// GetUint32 ...
func (c *RemoteConfig) GetUint32(key string) uint32 {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.cfg.GetUint32(key)
}

// GetUint64 ...
func (c *RemoteConfig) GetUint64(key string) uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.cfg.GetUint64(key)
}

// GetSizeInBytes ...
func (c *RemoteConfig) GetSizeInBytes(key string) uint {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.cfg.GetSizeInBytes(key)
}
