package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/spf13/cast"
)

// SyncMapKVStore ...
type SyncMapKVStore struct {
	kv *sync.Map
}

// NewSyncMapConfig returns KVStore based on sync.Map storage.
func NewSyncMapConfig(sm *sync.Map) *SyncMapKVStore {
	return &SyncMapKVStore{kv: sm}
}

// Set ...
func (cfg *SyncMapKVStore) Set(key string, value interface{}) {
	cfg.kv.Store(key, value)
}

// IsSet ...
func (cfg *SyncMapKVStore) IsSet(key string) bool {
	_, ok := cfg.kv.Load(key)

	return ok
}

// Get ...
func (cfg *SyncMapKVStore) Get(key string) interface{} {
	val, _ := cfg.kv.Load(key)

	return val
}

// AllKeys ...
func (cfg *SyncMapKVStore) AllKeys() []string {
	keys := []string{}

	cfg.kv.Range(func(key, _value interface{}) bool {
		k := fmt.Sprintf("%v", key)

		keys = append(keys, k)

		return true
	})

	return keys
}

// GetBool ...
func (cfg *SyncMapKVStore) GetBool(key string) bool {
	return cast.ToBool(cfg.Get(key))
}

// GetTime ...
func (cfg *SyncMapKVStore) GetTime(key string) (t time.Time) {
	return cast.ToTime(cfg.Get(key))
}

// GetDuration ...
func (cfg *SyncMapKVStore) GetDuration(key string) time.Duration {
	return cast.ToDuration(cfg.Get(key))
}

// GetInt ...
func (cfg *SyncMapKVStore) GetInt(key string) int {
	return cast.ToInt(cfg.Get(key))
}

// GetInt32 ...
func (cfg *SyncMapKVStore) GetInt32(key string) int32 {
	return cast.ToInt32(cfg.Get(key))
}

// GetInt64 ...
func (cfg *SyncMapKVStore) GetInt64(key string) int64 {
	return cast.ToInt64(cfg.Get(key))
}

// GetIntSlice ...
func (cfg *SyncMapKVStore) GetIntSlice(key string) []int {
	return cast.ToIntSlice(cfg.Get(key))
}

// GetString ...
func (cfg *SyncMapKVStore) GetString(key string) string {
	return cast.ToString(cfg.Get(key))
}

// GetStringSlice ...
func (cfg *SyncMapKVStore) GetStringSlice(key string) []string {
	return cast.ToStringSlice(cfg.Get(key))
}

// GetStringMap ...
func (cfg *SyncMapKVStore) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(cfg.Get(key))
}

// GetStringMapStringSlice ...
func (cfg *SyncMapKVStore) GetStringMapStringSlice(key string) map[string][]string {
	return cast.ToStringMapStringSlice(cfg.Get(key))
}

// GetFloat64 ...
func (cfg *SyncMapKVStore) GetFloat64(key string) float64 {
	return cast.ToFloat64(cfg.Get(key))
}

// GetUint ...
func (cfg *SyncMapKVStore) GetUint(key string) uint {
	return cast.ToUint(cfg.Get(key))
}

// GetUint32 ...
func (cfg *SyncMapKVStore) GetUint32(key string) uint32 {
	return cast.ToUint32(cfg.Get(key))
}

// GetUint64 ...
func (cfg *SyncMapKVStore) GetUint64(key string) uint64 {
	return cast.ToUint64(cfg.Get(key))
}

// GetSizeInBytes ...
func (cfg *SyncMapKVStore) GetSizeInBytes(key string) uint {
	sizeStr := cast.ToString(cfg.Get(key))
	return parseSizeInBytes(sizeStr)
}
