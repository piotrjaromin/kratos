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

var requiredAttackFlags = []string{"test-file", "duration", "users", "ramp-up"}

// StartServer creates cli command for creating server which can be run either in slave or master node
func StartServer(logger *log.Logger, conf *config.Config) cli.Command {
	return cli.Command{
		Name:    "attack",
		Aliases: []string{"a"},
		Usage:   "Commands run test from localhost",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "test-file",
				Usage: "File with defined test",
			},
			cli.StringFlag{
				Name:  "mode",
				Usage: fmt.Sprintf("Mode in which server should be started, validModes are %s", strings.Join(validModes, ", ")),
			},
		},
		Action: func(c *cli.Context) error {
			if err := validateRequiredFlags(c, requiredAttackFlags); err != nil {
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
