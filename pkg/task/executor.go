package task

type Executor struct {
	RegisterTasks  []DataTask
	ExecutionQueue chan DataTask
	ActiveTasks    map[string]DataTask
}

func (e *Executor) GetRegisteredTask() []DataTask {
	return e.RegisterTasks
}

func (e *Executor) AddActiveTask(dt DataTask) {
	e.ActiveTasks[dt.GetRunId()] = dt
}

func (e *Executor) ExecLoop() {
	for t := range e.ExecutionQueue {
		e.AddActiveTask(t)
		go t.Start()
	}
}

func (e *Executor) Start() {
	go func() {
		e.ExecLoop()
	}()
}

func NewExecutor(dataTasks []DataTask) *Executor {
	return &Executor{
		RegisterTasks:  dataTasks,
		ExecutionQueue: make(chan DataTask, 100),
		ActiveTasks:    make(map[string]DataTask),
	}
}
