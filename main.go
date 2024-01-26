package main

import (
	"log/slog"
	"os"

	"starling/cmd/server"
)

func main() {
	if err := run(); err != nil {
		slog.Error("Failed to start the app", "error", err)
	}
}

func run() error {
	options := server.Options{
		Logger: getLogger(),
	}
	return server.Run(&options)
}

func getLogger() *slog.Logger {
	if os.Getenv("DEBUG") != "true" {
		return slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}
	return slog.Default()
}
