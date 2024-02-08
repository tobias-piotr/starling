package worker

import (
	"fmt"
	"log/slog"

	"starling/internal/events"
)

type (
	Task  func() error
	Tasks map[string][]Task
)

type Worker struct {
	tasks    Tasks
	eventBus events.EventBus
}

func NewWorker(eventBus events.EventBus) *Worker {
	return &Worker{tasks: Tasks{}, eventBus: eventBus}
}

func (w Worker) AddTask(event string, task Task) {
	w.tasks[event] = append(w.tasks[event], task)
}

func (w Worker) Run() error {
	listener := make(chan map[string]any)

	go w.eventBus.Read(listener)

	for event := range listener {
		slog.Info("Received event", "event", event)

		// Get type
		typ, ok := event["type"]
		if !ok {
			err := w.backOff(event, []error{fmt.Errorf("event has no type: %v", event)})
			if err != nil {
				return fmt.Errorf("failed to back off: %w", err)
			}
			continue
		}
		typStr, ok := typ.(string)
		if !ok {
			err := w.backOff(event, []error{fmt.Errorf("event type is not a string: %v", event)})
			if err != nil {
				return fmt.Errorf("failed to back off: %w", err)
			}
			continue
		}

		// Execute
		errs := w.execute(typStr)
		var err error
		if len(errs) == 0 {
			err = w.confirm(event)
		} else {
			err = w.backOff(event, errs)
		}
		if err != nil {
			return fmt.Errorf("failed to confirm or back off: %w", err)
		}
	}

	return nil
}

// TODO: Check panics
// TODO: Handle retry
func (w Worker) execute(event string) []error {
	tasks, ok := w.tasks[event]
	if !ok {
		return nil
	}

	var errs []error
	for _, task := range tasks {
		if err := task(); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func (w Worker) confirm(event map[string]any) error {
	slog.Info("Event processed", "event_id", event["_id"], "event", event["type"])
	// TODO: Maybe pass entire event
	return w.eventBus.Confirm(event["_id"].(string))
}

func (w Worker) backOff(event map[string]any, errs []error) error {
	slog.Error("Event processing failed", "event_id", event["_id"], "event", event["type"], "errors", errs)
	return w.eventBus.BackOff(event)
}
