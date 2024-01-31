/*
Copyright © 2024 Piotr Tobiasz
*/
package server

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"starling/cmd"
	"starling/internal/database"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
)

var port string

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts server",
	Run: func(_ *cobra.Command, _ []string) {
		db, err := database.GetConnection(os.Getenv("DATABASE_DSN"))
		if err != nil {
			slog.Error("Failed to connect to database", "error", err)
			os.Exit(1)
		}
		if err := database.Migrate(db); err != nil {
			slog.Error("Failed to migrate database", "error", err)
			os.Exit(1)
		}

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

		e.Start(":" + port)
	},
}

func init() {
	cmd.RootCmd.AddCommand(serverCmd)

	// Server port
	serverCmd.Flags().StringVarP(&port, "port", "p", "8888", "Port to listen on")
}
