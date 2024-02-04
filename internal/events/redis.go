package events

type RedisEventBus struct{}

func NewRedisEventBus() *RedisEventBus {
	return &RedisEventBus{}
}

func (r *RedisEventBus) Publish(event Event) error {
	return nil
}
