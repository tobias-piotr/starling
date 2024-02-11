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
		red := getRedis()
		bus := events.NewRedisEventBus(
			red,
			&events.RedisBusArgs{
				Stream:        os.Getenv("REDIS_STREAM"),
				FailureStream: os.Getenv("REDIS_FAILURE_STREAM"),
				ConsumerGroup: os.Getenv("REDIS_CGROUP"),
				ConsumerName:  os.Getenv("REDIS_CNAME"),
			},
		)

		w := worker.NewWorker(bus)
		w.AddTask(trips.TripCreated{}.String(), func(_ map[string]any) error {
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
