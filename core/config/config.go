package config

import (
	"regexp"
	"text/template"
	"time"
)

type Key string

// ConfigStore is an interface on how a basic storage of
// config key-value pairs should have.
type ConfigStore interface {
	// SetEnvPrefix sets the prefix for env vars that will be treated as config,
	// it will only take affect before config has loaded.
	SetEnvPrefix(prefix string)
	// GetEnvPrefix gets the prefix setting.
	GetEnvPrefix() string
	// AddPath adds paths for config engine to read.
	AddPath(paths ...string)
	// Load loads configs from paths and current env.
	//
	// Load could be called multiple times, but recalling will not change the value
	// once a higher precedence level of the key has been set.
	Load() error
	// Dump exports the current effective config.
	Dump() map[string]string

	// SetDefault sets the default value for a given key.
	SetDefault(key Key, value any)
	// Set sets the value for a given key,
	// it will overwrite existing value if any.
	Set(key Key, value any)
	// Get gets the value for a given key.
	Get(key Key) any
}

// Config is an interface which extends the [ConfigStore]
// with common Get helpers
type Config interface {
	ConfigStore
	GetString(key Key) string
	GetBool(key Key) bool
	GetInt(key Key) int
	GetUint(key Key) uint32
	GetFloat(key Key) float32
	GetDuration(key Key) time.Duration
	GetTime(key Key) time.Time
	GetRegex(key Key) *regexp.Regexp
	GetTmpl(key Key) *template.Template
}
