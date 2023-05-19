package main

import (
	"os"

	"app/core/config"

	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slog"
)

var configFlag = &cli.StringFlag{
	Name:    "config",
	Usage:   "set config file path.",
	Aliases: []string{"f"},
	Action: func(ctx *cli.Context, s string) error {
		slog.Debug("config path added", slog.String("path", s))
		config.AddPath(s)
		return nil
	},
}

var workingDir = &cli.StringFlag{
	Name:  "cd",
	Usage: "change working directory to specific path",
	Action: func(ctx *cli.Context, s string) error {
		slog.Debug("process changed working dir", slog.String("dir", s))
		return os.Chdir(s)
	},
}
