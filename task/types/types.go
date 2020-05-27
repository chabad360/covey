package types

import "bytes"

// TaskPlugin defines the interface for Task module plugins.
type TaskPlugin interface {
	// NewTask creates a Task object and runs it in a go routine, and returns the Task.
	NewTask(taskJSON []byte) (ITask, error)
}

// ITask Defines the read methods for a task.
type ITask interface {
	// GetLog returns the full log of the task.
	GetLog() []string

	// GetUnread returns all unread log lines of the task.
	GetUnread() []string

	// GetID returns the ID of the task.
	GetID() [32]byte

	// GetIDShort returns the first 16 bytes of the task ID.
	GetIDShort() [8]byte
}

// Task defines the information of a task.
type Task struct {
	State   int           `json:"state,omitempty"`
	Plugin  string        `json:"plugin,omitempty"`
	ID      [32]byte      `json:"id,omitempty"`
	Node    string        `json:"node,omitempty"`
	Details interface{}   `json:"details,omitempty"`
	Log     []string      `json:"log,omitempty"`
	Buffer  *bytes.Buffer `json:"-"`
}

// GetID returns the ID of the task.
func (t *Task) GetID() [32]byte { return t.ID }

// GetIDShort returns the first 8 bytes of the task ID.
func (t *Task) GetIDShort() [8]byte { var x [8]byte; copy(x[:], t.ID[:8]); return x }
