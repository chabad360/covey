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
	sync.Mutex
	list     []task
	nonEmpty chan bool
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (q *queue) UnmarshalJSON(b []byte) error {
	var m map[int]task
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}

	q.Lock()
	for i := 0; i < len(m); i++ {
		q.list = append(q.list, m[i])
	}
	q.Unlock()

	q.full()
	return nil
}

// Get retrieves the first item in the queue. If the queue is empty, Get blocks until its not.
func (q *queue) Get() task {
	<-q.nE()

	q.Lock()
	var t task
	t, q.list = q.list[0], q.list[1:]
	q.Unlock()
	q.full()

	return t
}

func (q *queue) nE() <-chan bool {
	q.Lock()
	if q.nonEmpty == nil {
		q.nonEmpty = make(chan bool, 1)
	}
	d := q.nonEmpty
	q.Unlock()
	return d
}

func (q *queue) full() {
	q.nE()
	q.Lock()
	if len(q.list) > 0 && len(q.nonEmpty) == 0 {
		q.nonEmpty <- true
	}
	q.Unlock()
}

type runningTask struct {
	task
	log      []string
	ExitCode int
	State    int
	done     chan struct{}
	mu       *sync.Mutex
}

// GetLog returns the current log output
func (r *runningTask) GetLog() []string {
	var d []string
	r.mu.Lock()
	d, r.log = r.log, []string{}
	r.mu.Unlock()
	return d
}

// Log adds a line to the log.
func (r *runningTask) Log(log string) {
	r.mu.Lock()
	r.log = append(r.log, log)
	r.mu.Unlock()
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
		mu: &sync.Mutex{},
	}
	rt.task = t

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
