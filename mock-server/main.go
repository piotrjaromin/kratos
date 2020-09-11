package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", createHandler())
	e.Logger.Fatal(e.Start(":1323"))
}

func createHandler() echo.HandlerFunc {
	var requestCount uint64
	hello := func(c echo.Context) error {
		atomic.AddUint64(&requestCount, 1)
		return c.String(http.StatusOK, "Hello, World!")
	}

	counterRawInterval := 5
	ticker := time.NewTicker(time.Duration(counterRawInterval) * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				rps := requestCount / uint64(counterRawInterval)
				fmt.Printf("Rps %+v\n", rps)
				atomic.StoreUint64(&requestCount, 0)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	return hello
}
