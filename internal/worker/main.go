package worker

import (
	"fmt"
	"log/slog"

	"starling/internal/events"
)

type (
	Task  func(data map[string]any) error
	Tasks map[string][]Task
)

// Worker allows to specify tasks for specific event types
// It then uses the event bus to listen for events and execute tasks
type Worker struct {
	tasks    Tasks
	eventBus events.EventBus
}

func NewWorker(eventBus events.EventBus) *Worker {
	return &Worker{tasks: Tasks{}, eventBus: eventBus}
}

// AddTask adds a task to the map of tasks
func (w Worker) AddTask(event string, task Task) {
	w.tasks[event] = append(w.tasks[event], task)
}

// Run starts the worker, listens for events and executes tasks
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
				return err
			}
			continue
		}
		typStr, ok := typ.(string)
		if !ok {
			err := w.backOff(event, []error{fmt.Errorf("event type is not a string: %v", event)})
			if err != nil {
				return err
			}
			continue
		}

		// Execute
		errs := w.executeTasks(typStr, event)
		var err error
		if len(errs) == 0 {
			err = w.confirm(event)
		} else {
			err = w.backOff(event, errs)
		}
		if err != nil {
			return fmt.Errorf("execute confirm or back off: %w", err)
		}
	}

	return nil
}

// TODO: Check panics
// TODO: Handle retry
// executeTasks executes all tasks for the given event type
func (w Worker) executeTasks(eventType string, eventData map[string]any) []error {
	tasks, ok := w.tasks[eventType]
	if !ok {
		return nil
	}

	var errs []error
	for _, task := range tasks {
		if err := task(eventData); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

// confirm uses the event bus to confirm the event
func (w Worker) confirm(event map[string]any) error {
	slog.Info("Confirming event", "event_id", event["_id"], "event", event["type"])

	err := w.eventBus.Confirm(event["_id"].(string))
	if err != nil {
		return fmt.Errorf("confirm: %w", err)
	}

	return nil
}

// backOff uses the event bus to back off the event
func (w Worker) backOff(event map[string]any, errs []error) error {
	slog.Error("Backing off event", "event_id", event["_id"], "event", event, "errors", errs)

	err := w.eventBus.BackOff(event)
	if err != nil {
		return fmt.Errorf("back off: %w", err)
	}

	return nil
}
