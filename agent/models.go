package main

import (
	"fmt"
	json "github.com/json-iterator/go"
	"sync"
)

var closedChan = make(chan struct{})

func init() {
	close(closedChan)
}

// Queue is a basic concurrency-safe slice with extra checking to ensure that we don't pop an empty list.
type Queue struct {
	mutex    *sync.Mutex
	list     []task
	nonEmpty chan struct{}
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (q *Queue) UnmarshalJSON(b []byte) error {
	var m map[int]task
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}

	q.mutex.Lock()
	for i := 0; i < len(m); i++ {
		q.list = append(q.list, m[i])
	}
	q.mutex.Unlock()

	q.checkIfEmpty()
	return nil
}

// GetNext retrieves the pops the next item in the Queue. If the Queue is empty, GetNext blocks until its not.
func (q *Queue) GetNext() task {
	<-q.notEmpty()

	q.mutex.Lock()
	var t task
	t, q.list = q.list[0], q.list[1:]
	q.mutex.Unlock()
	q.checkIfEmpty()

	return t
}

// notEmpty returns a buffered channel that gets filled every time there is a new item added to the queue.
func (q *Queue) notEmpty() <-chan struct{} {
	q.mutex.Lock()
	if q.nonEmpty == nil {
		q.nonEmpty = make(chan struct{}, 1)
	}
	q.mutex.Unlock()
	return q.nonEmpty
}

// checkIfEmpty fills the nonEmpty channel if there is something in the queue.
func (q *Queue) checkIfEmpty() {
	q.notEmpty()
	q.mutex.Lock()
	if len(q.list) > 0 && len(q.nonEmpty) == 0 {
		q.nonEmpty <- struct{}{}
	}
	q.mutex.Unlock()
}

func newRunningTask(t task) runningTask {
	return runningTask{
		task:     t,
		context:  &baseContext{},
		mutex:    &sync.Mutex{},
		ExitCode: make(chan int, 1),
		State:    make(chan int, 1),
		log:      &[]string{},
	}
}

type runningTask struct {
	task
	ExitCode chan int
	State    chan int
	context  *baseContext
	mutex    *sync.Mutex
	log      *[]string
}

// GetLog returns the current log output.
func (r *runningTask) GetLog() (log []string) {
	r.mutex.Lock()
	fmt.Println(*r.log)
	log, *r.log = *r.log, []string{}
	r.mutex.Unlock()
	fmt.Println(log)
	return
}

// Log adds a line to the log.
func (r *runningTask) Log(log string) {
	r.mutex.Lock()
	*r.log = append(*r.log, log)
	r.mutex.Unlock()
}

// Finish marks the task as completed.
func (r *runningTask) Finish(exitCode int, state int) {
	r.mutex.Lock()

	r.ExitCode <- exitCode
	r.State <- state

	r.context.Close()
	r.mutex.Unlock()
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

// baseContext is the base implementation of a context.Context (doesn't fit the interface tho).
type baseContext struct {
	init  sync.Once
	mutex *sync.Mutex
	done  chan struct{}
}

func (c *baseContext) initMutex() {
	c.init.Do(func() { c.mutex = &sync.Mutex{} })
}

// Done returns a channel that will be closed once Close is run.
func (c *baseContext) Done() <-chan struct{} {
	c.initMutex()
	c.mutex.Lock()
	if c.done == nil {
		c.done = make(chan struct{})
	}
	c.mutex.Unlock()
	return c.done
}

// Close closes the channel returned by Done.
func (c *baseContext) Close() {
	c.initMutex()
	c.mutex.Lock()
	if c.done == nil {
		c.done = closedChan
	} else {
		close(c.done)
	}
	c.mutex.Unlock()
}
