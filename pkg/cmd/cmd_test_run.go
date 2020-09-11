package cmd

import (
	"fmt"
	"time"

	"github.com/piotrjaromin/kratos/pkg/attack"
	"github.com/piotrjaromin/kratos/pkg/config"
	"github.com/piotrjaromin/kratos/pkg/log"
	"github.com/piotrjaromin/kratos/pkg/utils"
	"github.com/urfave/cli"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

// TestRun creates cli command for running single test from localhost
func TestRun(logger *log.Logger, conf *config.Config) cli.Command {
	return cli.Command{
		Name:    "attack",
		Aliases: []string{"a"},
		Usage:   "Run single test from localhost",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:     "test-file",
				Usage:    "File with defined test",
				Required: true,
			},
			cli.StringFlag{
				Name:     "report-format",
				Usage:    "output format for runned test",
				Required: false,
				Value:    "img",
			},
			cli.IntFlag{
				Name:     "duration",
				Usage:    fmt.Sprintf("How long test should be run (without ramp-up-time) [seconds]"),
				Required: true,
			},
			cli.IntFlag{
				Name:     "ramp-up-time",
				Usage:    fmt.Sprintf("how long before reaching max rps value [seconds]"),
				Required: true,
			},
			cli.IntFlag{
				Name:     "max-rps",
				Usage:    fmt.Sprintf("How many rps test should reach"),
				Required: true,
			},
			cli.BoolFlag{
				Name:  "keep-alive",
				Usage: fmt.Sprintf("If keep-alive should be used"),
			},
			// TODO add common options file
		},
		Action: func(c *cli.Context) error {
			testFile := c.String("test-file")
			testContents, err := attack.TestFileProvider(testFile)
			if err != nil {
				return fmt.Errorf("Unable to read test file. Details: %s", err.Error())
			}

			duration := time.Second * time.Duration(c.Int("duration"))
			keepAlive := c.Bool("keep-alive")

			rampupTime := time.Second * time.Duration(c.Int("ramp-up-time"))
			rps := c.Int("max-rps")

			testDuration := rampupTime + duration
			logger.Infof("Rps are %+v, test duration is %+v, ramp up time is %+v\n", rps, testDuration, rampupTime)
			rate := attack.NewRampUpPacer(logger, rampupTime, rps)

			opts := attack.DefaultOpts(log.GetCLILogger(), testFile, testDuration, rate)
			opts.Keepalive = keepAlive

			out := &utils.StringWriterReader{}

			atk := vegeta.NewAttacker(
				vegeta.Timeout(opts.Timeout),
				// vegeta.TLSConfig(tlsc),
				// vegeta.KeepAlive(opts.Keepalive),
				vegeta.Workers(opts.Workers),
				vegeta.MaxWorkers(opts.MaxWorkers),
				vegeta.Connections(opts.Connections),
				vegeta.MaxConnections(opts.MaxConnections),
				vegeta.HTTP2(opts.HTTP2),
				vegeta.MaxBody(opts.MaxBody),
			)

			tr := attack.CreateTargeter(opts.Logger, testContents)
			if err := attack.Run(out, atk, tr, opts); err != nil {
				return fmt.Errorf("Error during test run. Details: %s", err.Error())
			}

			reportFormat := c.String("report-format")
			bucketStr := "" // If this will be used it cannot be empty
			if err := attack.Report(out, reportFormat, bucketStr, 10); err != nil {
				return fmt.Errorf("Error during test report generation. Details: %s", err.Error())
			}

			return nil
		},
	}
}
