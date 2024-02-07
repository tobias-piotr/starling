package worker

import (
	"log/slog"
)

type (
	Task  func() error
	Tasks map[string][]Task
)

type Worker struct {
	tasks Tasks
	// TODO: Add redis here and implement ack method
}

func NewWorker() *Worker {
	return &Worker{tasks: Tasks{}}
}

func (w *Worker) AddTask(event string, task Task) {
	w.tasks[event] = append(w.tasks[event], task)
}

// TODO: Check panics
func (w Worker) Execute(event string) []error {
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

// TODO: How to let redis know that we are done with the event?
func (w *Worker) Run(listener chan map[string]any) error {
	for event := range listener {
		slog.Info("Received event", "event", event)

		typ, ok := event["type"]
		if !ok {
			slog.Error("Event has no type", "event", event)
			continue
		}

		// TODO: Goroutine it
		w.Execute(typ.(string))
	}
	return nil
}
