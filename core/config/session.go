package config

import (
	"fmt"
	"regexp"
	"strconv"
	"text/template"
	"time"
)

// Session is an implementation of [Config], and offers common Get helpers
type Session struct {
	ConfigStore
}

// NewSession starts a new config session and uses [ConfigStore] as storage.
func NewSession(s ConfigStore) Config {
	return &Session{s}
}

func (s *Session) GetString(key Key) string {
	val := s.Get(key)
	if val == nil {
		return ""
	}
	return any2str(val)
}

func (s *Session) GetBool(key Key) bool {
	val := s.Get(key)
	if val != nil {
		switch v := val.(type) {
		case bool:
			return v
		case string:
			p, _ := strconv.ParseBool(v)
			return p
		}
	}
	return false
}

func (s *Session) GetInt(key Key) int {
	val := s.Get(key)
	if val != nil {
		switch v := val.(type) {
		case int:
			return v
		case string:
			p, _ := strconv.Atoi(v)
			return p
		}
	}
	return 0
}
func (s *Session) GetUint(key Key) uint32 {
	v := s.Get(key)
	if v != nil {
		switch val := v.(type) {
		case uint32:
			return val
		case string:
			p, _ := strconv.ParseUint(val, 10, 32)
			return uint32(p)
		}
	}
	return 0
}
func (s *Session) GetFloat(key Key) float32 {
	v := s.Get(key)
	if v != nil {
		switch val := v.(type) {
		case float32:
			return val
		case string:
			p, _ := strconv.ParseFloat(val, 32)
			return float32(p)
		}
	}
	return 0
}
func (s *Session) GetDuration(key Key) time.Duration {
	v := s.Get(key)
	if v != nil {
		switch val := v.(type) {
		case time.Duration:
			return val
		case string:
			p, _ := time.ParseDuration(val)
			return p
		}
	}
	return 0
}
func (s *Session) GetTime(key Key) time.Time {
	v := s.Get(key)
	if v != nil {
		switch val := v.(type) {
		case time.Time:
			return val
		case string:
			p, _ := time.Parse("2006-01-02", val)
			return p
		}
	}
	return time.Time{}
}
func (s *Session) GetRegex(key Key) *regexp.Regexp {
	v := s.Get(key)
	if v != nil {
		switch val := v.(type) {
		case *regexp.Regexp:
			return val
		case string:
			p, _ := regexp.Compile(val)
			return p
		}
	}
	return nil
}
func (s *Session) GetTmpl(key Key) *template.Template {
	v := s.Get(key)
	if v != nil {
		switch val := v.(type) {
		case *template.Template:
			return val
		case string:
			p, _ := template.New("").Parse(val)
			return p
		}
	}
	return nil
}

const (
	TRUE  = "true"
	FALSE = "false"
)

func any2str(val any) string {
	switch s := val.(type) {
	case string:
		return s
	case bool:
		if s {
			return TRUE
		}
		return FALSE
	case int:
		return strconv.Itoa(s)
	case uint:
		return strconv.FormatUint(uint64(s), 10)
	case uint32:
		return strconv.FormatUint(uint64(s), 10)
	case uint64:
		return strconv.FormatUint(s, 10)
	case float32:
		return strconv.FormatFloat(float64(s), byte('f'), -1, 32)
	case float64:
		return strconv.FormatFloat(s, byte('f'), -1, 64)
	case time.Duration:
		return s.String()
	case *regexp.Regexp:
		return s.String()
	default:
		return fmt.Sprint(val)
	}
}
