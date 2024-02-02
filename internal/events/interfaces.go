package events

type EventBus interface {
	Publish(event string, payload map[string]any) error
}
