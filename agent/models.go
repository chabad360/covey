package main

import (
	json "github.com/json-iterator/go"
	"sync"
)

var closedChan = make(chan struct{})

func init() {
	close(closedChan)
}

type queue struct {
	list     []task
	mu       *sync.Mutex
	nonEmpty chan bool
	init     bool
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (q *queue) UnmarshalJSON(b []byte) error {
	q.lazyInit()
	var m map[int]task
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}

	q.mu.Lock()
	for i := 0; i < len(m); i++ {
		q.list = append(q.list, m[i])
	}
	q.mu.Unlock()
	q.full()
	return nil
}

// Init initializes the queue.
func (q *queue) Init() {
	q.mu = &sync.Mutex{}
	q.mu.Lock()
	q.list = []task{}
	q.nonEmpty = make(chan bool, 1)
	q.init = true
	q.mu.Unlock()
}

// Get retrieves the first item in the queue. If the queue is empty, Get blocks until its not.
func (q *queue) Get() task {
	q.lazyInit()
	<-q.nonEmpty

	q.mu.Lock()
	var t task
	t, q.list = q.list[0], q.list[1:]
	q.mu.Unlock()
	q.full()

	return t
}

func (q *queue) lazyInit() {
	q.mu.Lock()
	if !q.init {
		q.mu.Unlock()
		q.Init()
	} else {
		q.mu.Unlock()
	}
}

func (q *queue) full() {
	q.lazyInit()
	q.mu.Lock()
	if len(q.list) > 0 && len(q.nonEmpty) == 0 {
		q.nonEmpty <- true
	}
	q.mu.Unlock()
}

type runningTask struct {
	Log      chan string
	ExitCode int
	State    int
	ID       string
	Cmd      string
	done     chan struct{}
	mu       *sync.Mutex
}

// Done returns a channel that will be closed when the task completes.
func (r *runningTask) Done() <-chan struct{} {
	r.mu.Lock()
	if r.done == nil {
		r.done = make(chan struct{})
	}
	d := r.done
	r.mu.Unlock()
	return d
}

// Finish marks the task as completed.
func (r *runningTask) Finish(exitCode int, state int) {
	r.mu.Lock()
	r.ExitCode = exitCode
	r.State = state
	if r.done == nil {
		r.done = closedChan
	} else {
		close(r.done)
	}
	r.mu.Unlock()
}

func newRunningTask(t task) *runningTask {
	rt := &runningTask{
		mu:  &sync.Mutex{},
		Log: make(chan string, 1024),
		ID:  t.ID,
		Cmd: t.Command,
	}

	return rt
}

type returnTask struct {
	Log      []string `json:"log,omitempty"`
	ExitCode int      `json:"exit_code"`
	State    int      `json:"state"`
	ID       string   `json:"id"`
}

type task struct {
	Command string `json:"command"`
	ID      string `json:"id"`
}

type config struct {
	AgentID   string
	LogLevel  string
	AgentPath string
}
