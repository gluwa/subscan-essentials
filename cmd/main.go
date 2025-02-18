package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/itering/subscan/configs"
	"github.com/itering/subscan/internal/observer"
	"github.com/itering/subscan/internal/script"
	"github.com/itering/subscan/internal/server/http"
	"github.com/itering/subscan/internal/service"
	"github.com/itering/subscan/util"
	"github.com/itering/substrate-api-rpc/websocket"
	"github.com/lmittmann/tint"
	"github.com/urfave/cli"
	"golang.org/x/exp/slog"
)

func envLogLevel() slog.Level {
	switch strings.TrimSpace(strings.ToLower(os.Getenv("GO_LOG"))) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func initSlog() {
	logger := slog.New(tint.NewHandler(os.Stderr, &tint.Options{Level: envLogLevel()}))
	slog.SetDefault(logger)
}

func main() {
	defer func() {
		websocket.Close()
	}()
	initSlog()
	if err := setupApp().Run(os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func setupApp() *cli.App {
	util.AddressType = "42"
	app := cli.NewApp()
	app.Name = "SUBSCAN"
	app.Usage = "SUBSCAN Backend Service, use -h get help"
	app.Version = "1.1"
	app.Action = func(*cli.Context) error { run(); return nil }
	app.Description = "SubScan Backend Service, substrate blockchain explorer"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "conf", Value: "../configs"},
	}
	app.Before = func(context *cli.Context) error {
		configs.Init()
		runtime.GOMAXPROCS(runtime.NumCPU())
		return nil
	}
	app.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "Start one worker, E.g substrate",
			Action: func(c *cli.Context) error {
				observer.Run(c.Args().Get(0))
				return nil
			},
		},
		{
			Name:  "install",
			Usage: "Create database and create default conf file",
			Action: func(c *cli.Context) error {
				script.Install(c.Parent().String("conf"))
				return nil
			},
		},
		{
			Name:  "CheckCompleteness",
			Usage: "Create blocks completeness",
			Action: func(c *cli.Context) error {
				script.CheckCompleteness()
				return nil
			},
		},
	}
	return app
}

func run() {
	svc := service.NewReadOnly()

	httpSrv := http.NewHTTPServer(configs.Boot.Server, svc)
	defer func() {
		// Micro services
		svc.Close()
	}()

	runIt := func() error {
		if addr := configs.Boot.Server.Http.Addr; addr != "" {
			return httpSrv.Run(addr)
		} else {
			return httpSrv.Run()
		}
	}
	if err := runIt(); err != nil {
		panic(err)
	}
}
