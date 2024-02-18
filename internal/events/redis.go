package events

import (
	"context"
	"fmt"
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

type RedisBusArgs struct {
	Stream        string
	FailureStream string
	ConsumerGroup string
	ConsumerName  string
}

// RedisEventBus is an implementation of EventBus that uses Redis as the message broker.
type RedisEventBus struct {
	client        *redis.Client
	stream        string
	failureStream string
	consumerGroup string
	consumerName  string
}

func NewRedisEventBus(client *redis.Client, args *RedisBusArgs) *RedisEventBus {
	return &RedisEventBus{client, args.Stream, args.FailureStream, args.ConsumerGroup, args.ConsumerName}
}

// Publish publishes the event to the stream.
func (b *RedisEventBus) Publish(event Event) error {
	slog.Info("Publishing event", "event", event.String())

	err := b.client.XAdd(context.Background(), &redis.XAddArgs{
		Stream: b.stream,
		Values: event.Payload(),
	}).Err()
	if err != nil {
		return fmt.Errorf("publish event: %w", err)
	}

	return nil
}

// Read reads events from the stream and sends them to the listener.
func (b *RedisEventBus) Read(listener chan map[string]any) error {
	// Create the consumer group
	// Its ok if it already exists
	err := b.client.XGroupCreate(context.Background(), b.stream, b.consumerGroup, "0").Err()
	if err != nil {
		slog.Warn("Error creating consumer group", "error", err)
	}

	for {
		data, err := b.
			client.
			XReadGroup(context.Background(), &redis.XReadGroupArgs{
				Group:    b.consumerGroup,
				Consumer: b.consumerName,
				Streams:  []string{b.stream, ">"},
				Count:    5,
			}).
			Result()
		if err != nil {
			return fmt.Errorf("read events: %w", err)
		}

		for _, result := range data {
			for _, message := range result.Messages {
				message.Values["_id"] = message.ID
				listener <- message.Values
			}
		}
	}
}

// Confirm acknowledges the event.
func (b *RedisEventBus) Confirm(id string) error {
	err := b.client.XAck(context.Background(), b.stream, b.consumerGroup, id).Err()
	if err != nil {
		return fmt.Errorf("ack event: %w", err)
	}
	return nil
}

// BackOff moves the event to the failure stream and acknowledges it from the main stream.
func (b *RedisEventBus) BackOff(event map[string]any) error {
	// Add the event to the failure stream
	err := b.client.XAdd(context.Background(), &redis.XAddArgs{
		Stream: b.failureStream,
		Values: event,
	}).Err()
	if err != nil {
		return fmt.Errorf("add event to failure stream: %w", err)
	}

	// Ack the event from the main stream
	err = b.client.
		XAck(context.Background(), b.stream, b.consumerGroup, event["_id"].(string)).
		Err()
	if err != nil {
		return fmt.Errorf("ack event from main stream: %w", err)
	}

	return nil
}
