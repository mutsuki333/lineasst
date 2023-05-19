/*
	msg.go

	Purpose: To translate messages to user's perfered language.

	@author Evan Chen

	MODIFICATION HISTORY
	   Date        Ver    Name     Description
	---------- ------- ----------- -------------------------------------------
	2023/03/02  v1.0.0 Evan Chen   Initial release
*/

// Package msg is a package to generate translated messages
package msg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"

	"app/core/property"

	"google.golang.org/grpc/metadata"
)

// Pack is the file structure for a message collection
type Pack struct {
	Messages []impMsg `json:"messages"`
	Locale   string   `json:"locale"`
}

type impMsg struct {
	Key    string `json:"key"`
	Locale string `json:"locale"`
	Tmpl   string `json:"tmpl"`
}

// store is the collection of current registered messages
var store map[string]map[string]*template.Template

func init() {
	store = make(map[string]map[string]*template.Template)
}

// T is a function to translate messages, and any occurred errors will only be logged.
//
//   - It will use the default locale if the given lang is not found.
//   - If the message key is not found in the default locale, it will return the message key as-is with the data appended at the end.
func T(key, locale string, data any) string {

	if _, has_loc := store[locale]; has_loc {
		tmpl, ok := store[locale][key]
		if ok {
			return execute(tmpl, data)
		}
	}
	tmpl, ok := store[property.DefaultLocale][key]
	if ok {
		return execute(tmpl, data)
	}

	return fmt.Sprintf("%s: %+v", key, data)
}

// Tctx is like [T] but it uses the locale in the context.
//
// Note that it looks for the "x-accept-language" key
// which should be injected by the middleware.
func Tctx(ctx context.Context, key string, data any) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return T(key, property.DefaultLocale, data)
	}
	if r := md.Get("x-accept-language"); len(r) > 0 {
		return T(key, r[0], data)
	}
	return T(key, property.DefaultLocale, data)
}

type plain struct {
	Info string
}

// Plain returns a structure with {{.Info}} set to string.
func Plain(msg string) plain {
	return plain{Info: msg}
}

func execute(tmpl *template.Template, payload any) string {
	buff := &bytes.Buffer{}
	if err := tmpl.Execute(buff, payload); err != nil {
		return err.Error()
	}
	return buff.String()
}

// Load loads a message pack and adds to the current pack,
// it will overwrite existing messages.
func Load(pk *Pack) error {
	pklocal := pk.Locale
	if pklocal == "" {
		pklocal = property.DefaultLocale
	}
	store_locale(pklocal)

	for _, msg := range pk.Messages {
		tmpl, err := template.New("").Parse(msg.Tmpl)
		if err != nil {
			return err
		}
		if msg.Locale != "" {
			store_locale(msg.Locale)
			store[msg.Locale][msg.Key] = tmpl
		} else {
			store[pklocal][msg.Key] = tmpl
		}
	}

	return nil
}

// store_locale make sures a given locale is initiated
func store_locale(locale string) {
	if exist := store[locale]; exist == nil {
		store[locale] = make(map[string]*template.Template)
	}
}

// LoadFS is like Load, but it loads message packs from the given fs.FS
func LoadFS(f fs.FS, path string) error {

	stat, err := fs.Stat(f, path)
	if err != nil {
		return err
	}
	if stat.IsDir() {
		return fs.WalkDir(f, path, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			if filepath.Ext(path) == ".json" {
				content, err := fs.ReadFile(f, path)
				if err != nil {
					return err
				}
				var pak Pack
				if err := json.Unmarshal(content, &pak); err != nil {
					return err
				}
				return Load(&pak)
			}

			return nil
		})

	} else {
		content, err := fs.ReadFile(f, path)
		if err != nil {
			return err
		}
		var pak Pack
		if err := json.Unmarshal(content, &pak); err != nil {
			return err
		}
		return Load(&pak)
	}
}

// LoadPath is like Load, but it loads message packs from the filesystem
func LoadPath(paths ...string) error {
	for _, path := range paths {
		if err := LoadFS(os.DirFS(path), "."); err != nil {
			return err
		}
	}
	return nil
}
