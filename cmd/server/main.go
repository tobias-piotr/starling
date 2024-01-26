package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Options struct {
	Logger *slog.Logger
}

func main() {
	if err := Run(nil); err != nil {
		slog.Error("Failed to run server", "error", err)
	}
}

func Run(opts *Options) error {
	e := echo.New()

	e.Use(NewLoggingMiddleware(opts))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	e.Use(middleware.RequestID())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{Timeout: time.Second * 60}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	return e.Start(":8888")
}
