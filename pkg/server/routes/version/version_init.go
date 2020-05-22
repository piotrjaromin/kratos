package version

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/piotrjaromin/kratos/pkg/config"
)

type versionResponse struct {
	ServiceName string `json:"serviceName"`
	Commit      string `json:"commit"`
	Version     string `json:"version"`
}

// Init version endpoint
func Init(e *echo.Echo, config *config.Config) {
	versionHandler := func(c echo.Context) error {
		return c.JSON(http.StatusOK, versionResponse{
			ServiceName: config.Name(),
			Commit:      config.GitCommit(),
			Version:     config.Version(),
		})
	}

	e.GET("/version", versionHandler)
}
