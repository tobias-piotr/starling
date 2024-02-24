/*
Copyright © 2024 Piotr Tobiasz
*/
package cmd

import (
	"log/slog"
	"os"

	"starling/cmd"
	"starling/internal/ai"
	"starling/internal/database"
	"starling/internal/events"
	"starling/internal/worker"
	"starling/trips"
	tripsDB "starling/trips/database"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
)

func getDB() *sqlx.DB {
	db, err := database.GetConnection(os.Getenv("DATABASE_DSN"))
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
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

// workerCmd represents the worker command
var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Starts worker",
	Run: func(_ *cobra.Command, _ []string) {
		// Database
		db := getDB()
		tripsRepo := tripsDB.NewTripRepository(db)

		// Redis
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

		// OpenAI
		aiClient := ai.NewOpenAIClient(os.Getenv("OPENAI_KEY"))

		w := worker.NewWorker(bus)
		w.AddTask(trips.TripCreated{}.String(), func(data map[string]any) error {
			slog.Info("Trip created", "trip_id", data["trip_id"])
			return nil
		})
		w.AddTask(trips.TripRequested{}.String(), func(data map[string]any) error {
			slog.Info("Trip requested", "trip_id", data["trip_id"])
			return trips.RequestTrip(tripsRepo, aiClient, data["trip_id"].(string))
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
