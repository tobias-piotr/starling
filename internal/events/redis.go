package events

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(addr string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}

	return client, nil
}

type RedisEventBus struct {
	client *redis.Client
	stream string
}

func NewRedisEventBus(client *redis.Client, stream string) *RedisEventBus {
	return &RedisEventBus{client, stream}
}

func (b *RedisEventBus) Publish(event Event) error {
	slog.Info("Publishing event", "event", event.String())

	err := b.client.XAdd(context.Background(), &redis.XAddArgs{
		Stream: b.stream,
		Values: event.Payload(),
	}).Err()
	if err != nil {
		slog.Error("Failed to publish event", "error", err)
		return err
	}

	return nil
}

func (b *RedisEventBus) Read(listener chan map[string]any) error {
	data, err := b.client.XRead(context.Background(), &redis.XReadArgs{
		Streams: []string{b.stream, "0"},
		Count:   5,
		// TODO: Try with bigger, and without
		Block: 0,
	}).Result()
	if err != nil {
		slog.Error("Failed to read from stream", "error", err)
		return err
	}

	for _, result := range data {
		for _, message := range result.Messages {
			message.Values["_id"] = message.ID
			listener <- message.Values
		}
	}
	return nil
}
