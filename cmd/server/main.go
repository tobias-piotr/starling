package server

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Options struct {
	DB *sqlx.DB
}

func main() {
	if err := Run(nil); err != nil {
		slog.Error("Failed to run server", "error", err)
		os.Exit(1)
	}
}

func Run(opts *Options) error {
	e := echo.New()

	e.Use(NewLoggingMiddleware())
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
