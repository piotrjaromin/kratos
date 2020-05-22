package main

import (
	"os"

	"github.com/urfave/cli"

	"github.com/piotrjaromin/kratos/pkg/cmd"
	"github.com/piotrjaromin/kratos/pkg/config"
	"github.com/piotrjaromin/kratos/pkg/log"
)

// initialized by ldflags during build
var gitCommit = "n/a"

// Manually bumped
var version = "0.0.1"

func main() {
	appInfo := config.AppInfo{
		Version:   version,
		GitCommit: gitCommit,
	}

	conf, err := config.GetConfig(appInfo, nil)
	if err != nil {
		panic(err)
	}

	logger := log.GetLogger(conf)
	app := cli.NewApp()

	app.Commands = []cli.Command{
		cmd.StartServer(logger, conf),
	}

	app.Name = "Performance/load testing tool with built in server"
	app.Version = "0.0.1"
	app.Usage = ""

	err = app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
