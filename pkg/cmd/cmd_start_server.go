package cmd

import (
	"fmt"
	"strings"

	"github.com/piotrjaromin/kratos/pkg/config"
	"github.com/piotrjaromin/kratos/pkg/log"
	"github.com/piotrjaromin/kratos/pkg/server"
	"github.com/thoas/go-funk"
	"github.com/urfave/cli"
)

var requiredServerFlags = []string{"mode"}
var validModesWithDefaultPorts = map[server.Mode]string{
	server.Slave:  "8081",
	server.Master: "8080",
}
var validModes = []string{string(server.Slave), string(server.Master)}

// StartServer creates cli command for creating server which can be run either in slave or master node
func StartServer(logger *log.Logger, conf *config.Config) cli.Command {
	return cli.Command{
		Name:    "server",
		Aliases: []string{"s"},
		Usage:   "Starts server which can run performance tests",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "port",
				Usage: "Port on which server should be started, if not provided default values will be used",
			},
			cli.StringFlag{
				Name:  "mode",
				Usage: fmt.Sprintf("Mode in which server should be started, validModes are %s", strings.Join(validModes, ", ")),
			},
		},
		Action: func(c *cli.Context) error {
			if err := validateRequiredFlags(c, requiredServerFlags); err != nil {
				return err
			}

			port := c.String("port")
			mode := server.Mode(c.String("mode"))

			if !funk.Contains(validModesWithDefaultPorts, mode) {
				return fmt.Errorf("Invalid mode, must be one of %s", strings.Join(validModes, ", "))
			}

			if len(port) == 0 {
				port = validModesWithDefaultPorts[mode]
			}

			return server.Start(logger, conf, port, mode)
		},
	}
}
