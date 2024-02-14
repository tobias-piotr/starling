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
	"starling/internal/events"
	trips_api "starling/trips/api"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
)

var (
	port   string
	prefix string
)

func getDB() *sqlx.DB {
	db, err := database.GetConnection(os.Getenv("DATABASE_DSN"))
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	if err := database.Migrate(db); err != nil {
		slog.Error("Failed to migrate database", "error", err)
		os.Exit(1)
	}
	return db
}

func getRedis() *redis.Client {
	client, err := events.NewRedisClient(os.Getenv("REDIS_ADDR"))
	if err != nil {
		slog.Error("Failed to connect to redis", "error", err)
		os.Exit(1)
	}
	return client
}

func createServer() *echo.Echo {
	e := echo.New()

	e.Use(NewLoggingMiddleware())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	e.Use(middleware.RequestID())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{Timeout: time.Second * 60}))
	e.Pre(middleware.RemoveTrailingSlash())

	return e
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts server",
	Run: func(_ *cobra.Command, _ []string) {
		db := getDB()
		redisClient := getRedis()
		e := createServer()

		router := e.Group(prefix)
		router.GET("/health", func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
		})
		trips_api.Register(router, db, redisClient)

		e.Start(":" + port)
	},
}

func init() {
	cmd.RootCmd.AddCommand(serverCmd)

	// Server port
	serverCmd.Flags().StringVarP(&port, "port", "p", "8888", "Port to listen on")
	// Server prefix
	serverCmd.Flags().StringVarP(&prefix, "prefix", "r", "/sl", "Prefix for all routes")
}
