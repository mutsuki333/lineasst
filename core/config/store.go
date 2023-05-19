package config

import (
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

// Store is an implementation of ConfigStore
// that stores config values in an internal sync.Map .
//
// Store is ready to use without the need to initialize. Usage of store is concurrency safe.
type Store struct {
	paths  []string
	prefix string
	store  sync.Map
}

type setlvl int

const (
	lvl_default setlvl = iota
	lvl_file
	lvl_env
	lvl_set
)

type lvl_value struct {
	val any
	lvl setlvl
}

/*
	Functions to comply with Config interface
*/

func (s *Store) SetEnvPrefix(prefix string) {
	s.prefix = prefix
}
func (s *Store) GetEnvPrefix() string { return s.prefix }
func (s *Store) AddPath(paths ...string) {
	if s.paths == nil {
		s.paths = paths
	} else {
		s.paths = append(s.paths, paths...)
	}
}
func (s *Store) Load() error {
	envMap, _ := godotenv.Read(s.paths...)
	for k, v := range envMap {
		ret, ok := s.store.Load(Key(k))
		if !ok || ret.(lvl_value).lvl <= lvl_file {
			s.store.Store(Key(k), lvl_value{
				lvl: lvl_file,
				val: v,
			})
		}
	}
	if s.prefix == "" {
		s.prefix = "APP_"
	}
	envKVs := os.Environ()
	for _, envKV := range envKVs {
		if !strings.HasPrefix(envKV, s.prefix) {
			continue
		}
		envKV = strings.TrimPrefix(envKV, s.prefix)
		k, v, _ := strings.Cut(envKV, "=")
		ret, ok := s.store.Load(Key(k))
		if !ok || ret.(lvl_value).lvl <= lvl_env {
			s.store.Store(Key(k), lvl_value{
				lvl: lvl_env,
				val: v,
			})
		}
	}

	return nil
}

func (s *Store) Dump() map[string]string {
	dump := make(map[string]string)
	s.store.Range(func(key, value any) bool {
		dump[string(key.(Key))] = any2str(value.(lvl_value).val)
		return true
	})

	return dump
}

func (s *Store) SetDefault(key Key, value any) {
	ret, ok := s.store.Load(key)
	if !ok || ret.(lvl_value).lvl <= lvl_default {
		s.store.Store(key, lvl_value{lvl: lvl_default, val: value})
	}
}
func (s *Store) Set(key Key, value any) {
	s.store.Store(key, lvl_value{lvl: lvl_set, val: value})
}
func (s *Store) Get(key Key) any {
	v, ok := s.store.Load(key)
	if !ok {
		return nil
	}
	return v.(lvl_value).val
}
