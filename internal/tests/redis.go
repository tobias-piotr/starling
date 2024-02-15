package tests

import (
	"os"

	"starling/internal/events"
)

// GetRedisPublisher returns RedisEventBus that is connected to test stream
func GetRedisPublisher() *events.RedisEventBus {
	client, _ := events.NewRedisClient(os.Getenv("REDIS_ADDR"))
	return events.NewRedisEventBus(client, &events.RedisBusArgs{Stream: "tests"})
}
