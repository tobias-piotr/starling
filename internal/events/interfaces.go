package events

type Event interface {
	String() string
	Payload() map[string]any
}

type EventBus interface {
	Publish(event Event) error
	Read(listener chan map[string]any) error
}
