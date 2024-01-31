package worker

type (
	Task  func() error
	Tasks map[string][]Task
)

type Worker struct {
	tasks Tasks
}

func NewWorker() *Worker {
	return &Worker{tasks: Tasks{}}
}

func (w *Worker) AddTask(event string, task Task) {
	w.tasks[event] = append(w.tasks[event], task)
}

// TODO: Execute tasks with goroutines
// TODO: What if something panics?
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

func (w *Worker) Run() error {
	// TODO: Use signal to wait for events
	w.Execute("start")
	w.Execute("stop")
	return nil
}
