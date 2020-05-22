package server

import (
	"fmt"

	"github.com/piotrjaromin/kratos/pkg/log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/piotrjaromin/kratos/pkg/config"
	"github.com/piotrjaromin/kratos/pkg/server/routes/health"
	"github.com/piotrjaromin/kratos/pkg/server/routes/version"
)

var skipEndpoints = []string{"/version", "/health"}

// Mode of server in which it should be started
type Mode string

// Slave mode for server
const Slave Mode = "slave"

// Master mode for server
const Master Mode = "master"

// Start in server mode
func Start(logger *log.Logger, conf *config.Config, port string, mode Mode) error {
	e := getEcho(conf)

	logger.Info("Setting up routes")

	version.Init(e, conf)
	health.Init(e)

	logger.Infof("Starting server on port %s, in mode %s", port, mode)
	return e.Start(fmt.Sprintf(":%s", port))
}

func getEcho(conf *config.Config) *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: conf.Cors.AllowedOrigins,
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	return e
}
