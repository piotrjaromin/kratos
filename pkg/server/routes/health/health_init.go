package health

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type healthResponse struct {
	Message string `json:"message"`
}

func Init(e *echo.Echo) {
	versionHandler := func(c echo.Context) error {
		return c.JSON(http.StatusOK, healthResponse{
			Message: "OK",
		})
	}

	e.GET("/health", versionHandler)
}
