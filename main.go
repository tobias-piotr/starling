package main

import (
	"log/slog"
	"os"

	"starling/cmd/server"
	"starling/internal/database"
)

func main() {
	if err := run(); err != nil {
		slog.Error("Failed to start the app", "error", err)
		os.Exit(1)
	}
}

func run() error {
	db, err := database.GetConnection(os.Getenv("DATABASE_DSN"))
	if err != nil {
		return err
	}
	options := server.Options{DB: db}
	return server.Run(&options)
}

func setLogger() {
	if os.Getenv("DEBUG") != "true" {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	}
}
