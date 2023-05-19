package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	ossignal "os/signal"
	"syscall"
	"time"

	"app/core/property"
	"app/core/service"
	"app/core/signal"
	"app/core/util"

	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slog"
)

func main() {
	app := &cli.App{
		Name:    BinName,
		Usage:   Usage,
		Version: fmt.Sprintf("%s-%s %s", Version, ID, Build),
		Commands: []*cli.Command{
			StartCMD,
		},
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error(fmt.Sprintf("%s failed", BinName), util.ErrAtrr(err))
		os.Exit(1)
	}
}

var StartCMD = &cli.Command{
	Name:  "start",
	Usage: "start server",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "migrate",
			Usage: "auto migrate database.",
		},
		&cli.BoolFlag{
			Name:  "with-demo-data",
			Usage: "setup initial data for the system to work properly. This option will not override any existing data, and will skip on conflict.",
		},
		configFlag, workingDir,
	},
	Action: func(ctx *cli.Context) error {
		//-------------------------------------------------
		//- Load Configurations                           -
		//-------------------------------------------------
		if err := setup_config(); err != nil {
			return fmt.Errorf("setup error: %s", err)
		}

		//-------------------------------------------------
		//- Connection to database                        -
		//-------------------------------------------------
		// client, err := connect(ctx.Bool("migrate"))
		// if err != nil {
		// 	return fmt.Errorf("connection failed: %s", err)
		// }

		//-------------------------------------------------
		//- Setup Server                                  -
		//-------------------------------------------------
		srv, err := setup_service()
		if err != nil {
			return fmt.Errorf("server failed: %s", err)
		}

		//-------------------------------------------------
		//- Service Initiate and Load                     -
		//-------------------------------------------------
		if err := service.Initiate(); err != nil {
			return fmt.Errorf("lifecycle.Initiate failed: %s", err)
		}
		service.Load()
		signal.Handle(property.RELOAD, func(e signal.Event) { service.Load() })

		//-------------------------------------------------
		//- Application Start                             -
		//-------------------------------------------------
		property.SetState(property.STATE_PRESTART)
		srv.Start()

		// setup data
		// setAutoLogin(client)
		// if err = dbinit.Start(client, ctx.Bool("with-demo-data")); err != nil {
		// 	return err
		// }

		property.SetState(property.STATE_STARTED) // application started

		//-------------------------------------------------
		//- Graceful shutdown                             -
		//-------------------------------------------------
		quit := make(chan os.Signal, 1)
		ossignal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		term := <-quit // Application is blocked here
		slog.Info("signal received",
			slog.Any("sig", term),
			slog.String("act", "terminate"))
		property.SetState(property.STATE_TERM)
		srv.Stop(time.Second * 10)
		shutSvc := service.Del(time.Second * 10)
		<-shutSvc.Done()
		// client.Close()
		if errors.Is(shutSvc.Err(), context.DeadlineExceeded) {
			return fmt.Errorf("service termination timeout reached! force terminated")
		}

		return nil
	},
}
