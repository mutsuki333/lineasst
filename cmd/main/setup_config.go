package main

import (
	"fmt"
	"path/filepath"

	"app/core/auth"
	"app/core/config"
	"app/core/logger"
	"app/core/msg"
	"app/core/property"
	"app/core/server"
	"app/src/messages"

	"golang.org/x/exp/slog"
)

// setup_config loads configs and setup basic application settings.
func setup_config() error {

	if err := config.Load(); err != nil {
		return err
	}

	if config.GetString(property.DEBUG) == config.TRUE {
		property.RUN_MODE = property.DEBUG_MODE
	} else {
		property.RUN_MODE = property.PROD_MODE
	}

	if config.GetString(property.SECRET) != "" {
		auth.Secret = []byte(config.GetString(property.SECRET))
	}

	if property.IsDebug() {
		for k, v := range config.Dump() {
			fmt.Printf("%s=%s\n", k, v)
		}
	}

	//-------------------------------------------------
	//- Setup Logging                                 -
	//-------------------------------------------------

	property.SetLogLevel(config.GetString(property.LOG_LEVEL))

	if property.IsDebug() {
		slog.SetDefault(logger.NewLogger(BinName, &slog.HandlerOptions{AddSource: true}))
	} else {
		slog.SetDefault(logger.NewLogger(BinName))
	}
	server.SetLogger(logger.NewLogger("access"))

	//-------------------------------------------------
	//- Load language packs                           -
	//-------------------------------------------------

	if err := msg.LoadFS(messages.FS, "."); err != nil { // default message packs
		slog.Warn("load default message pack failed", slog.String("mod", "main"), slog.String("act", "setup"))
	}
	cust_msg := filepath.Join(config.GetString(property.CUSTOM), property.CUST_TMPL)
	if err := msg.LoadPath(cust_msg); err != nil { // overwrite existing
		slog.Warn("load custom message pack failed",
			slog.String("err", err.Error()),
			slog.String("mod", "main"),
			slog.String("act", "setup"),
			slog.String("path", cust_msg))
	}

	return nil
}
