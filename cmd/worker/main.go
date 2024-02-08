/*
Copyright © 2024 Piotr Tobiasz
*/
package cmd

import (
	"log/slog"
	"os"

	"starling/cmd"
	"starling/internal/events"
	"starling/internal/worker"
	"starling/trips"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
)

func getRedis() *redis.Client {
	client, err := events.NewRedisClient(os.Getenv("REDIS_ADDR"))
	if err != nil {
		slog.Error("Failed to connect to redis", "error", err)
		os.Exit(1)
	}
	return client
}

// workerCmd represents the worker command
var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Starts worker",
	Run: func(_ *cobra.Command, _ []string) {
		// TODO: Something needs to be done about stream name
		red := getRedis()
		bus := events.NewRedisEventBus(red, "trips", "trips-failures")

		w := worker.NewWorker(bus)
		w.AddTask(trips.TripCreated{}.String(), func() error {
			slog.Info("Trip created")
			return nil
		})

		if err := w.Run(); err != nil {
			slog.Error("Worker failed", "error", err)
			os.Exit(1)
		}
	},
}

func init() {
	cmd.RootCmd.AddCommand(workerCmd)
}
