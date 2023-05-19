/*
	default.go
	Purpose: Start a default config session and export all its functions
	so they could be used easier.

	@author Evan Chen
	@version 1.0 2023/02/22
*/

package config

import (
	"regexp"
	"text/template"
	"time"
)

var std Config

func init() {
	std = NewSession(&Store{})
}

// SetEnvPrefix sets the prefix for env vars that will be treated as config,
// it will only take affect before config has loaded.
func SetEnvPrefix(prefix string) { std.SetEnvPrefix(prefix) }

// GetEnvPrefix gets the prefix setting.
func GetEnvPrefix() string { return std.GetEnvPrefix() }

// AddPath adds paths for config engine to read.
func AddPath(paths ...string) { std.AddPath(paths...) }

// Load loads configs from paths and current env.
//
// Load could be called multiple times, but recalling will not change the value
// once a higher precedence level of the key has been set.
func Load() error { return std.Load() }

// Dump exports the current effective config.
func Dump() map[string]string { return std.Dump() }

// SetDefault sets the default value for a given key.
func SetDefault(key Key, value any) { std.SetDefault(key, value) }

// Set sets the value for a given key,
// it will overwrite existing value if any.
func Set(key Key, value any) { std.Set(key, value) }

func GetString(key Key) string           { return std.GetString(key) }
func GetBool(key Key) bool               { return std.GetBool(key) }
func GetInt(key Key) int                 { return std.GetInt(key) }
func GetUint(key Key) uint32             { return std.GetUint(key) }
func GetFloat(key Key) float32           { return std.GetFloat(key) }
func GetDuration(key Key) time.Duration  { return std.GetDuration(key) }
func GetTime(key Key) time.Time          { return std.GetTime(key) }
func GetRegex(key Key) *regexp.Regexp    { return std.GetRegex(key) }
func GetTmpl(key Key) *template.Template { return std.GetTmpl(key) }
