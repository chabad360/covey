package types

import (
	"container/list"
	"encoding/hex"
	"time"

	json "github.com/json-iterator/go"
)

// AgentTask contains the information about a task that is send to an agent.
type AgentTask struct {
	Command string `json:"command"`
	ID      string `json:"id"`
}

// TaskList implements List type as well as the json.Marshaler interface.
type TaskList struct{ list.List }

// TaskInfo contains new information about a running task.
type TaskInfo struct {
	Log      []string `json:"log"`
	ExitCode int      `json:"exit_code"`
	ID       string   `json:"id"`
}

// MarshalJSON implements the json.Marshaler interface.
func (l *TaskList) MarshalJSON() ([]byte, error) {
	m := make(map[int]AgentTask)
	i := 0
	for e := l.Front(); e != nil; e = e.Next() {
		m[i] = e.Value.(AgentTask)
		i++
	}
	return json.Marshal(m)
}

// TaskPlugin defines the interface for Task module plugins.
type TaskPlugin interface {
	// GetCommand returns the command to run the server.
	GetCommand(taskJSON []byte) (string, error)
}

// Task defines the information of a task.
type Task struct {
	State    TaskState         `json:"state"`
	Plugin   string            `json:"plugin"`
	ID       string            `json:"id"`
	Node     string            `json:"node"`
	Details  map[string]string `json:"details"`
	Log      []string          `json:"log"`
	Time     time.Time         `json:"time"`
	ExitCode int               `json:"exit_code"`
	Command  string            `json:"-"`
}

// GetID returns the ID of the task.
func (t *Task) GetID() string { return t.ID }

// GetIDShort returns the first 8 bytes of the task ID.
func (t *Task) GetIDShort() string { x, _ := hex.DecodeString(t.ID); return hex.EncodeToString(x[:8]) }

// GetState returns the current state of the task.
func (t *Task) GetState() TaskState { return t.State }

// GetPlugin returns the plugin of the task.
func (t *Task) GetPlugin() string { return t.Plugin }

// GetNode returns the node of the task.
func (t *Task) GetNode() string { return t.Node }

// GetTime returns the time of the task.
func (t *Task) GetTime() time.Time { return t.Time }

// GetDetails returns the details of the task.
func (t *Task) GetDetails() map[string]string { return t.Details }

// GetExitCode returns the exit code of the task.
func (t *Task) GetExitCode() int { return t.ExitCode }

// GetLog reads the unread buffer and adds it to the task's log, then returns that log.
func (t *Task) GetLog() []string { return t.Log }
