package types

import (
	"bytes"
	"encoding/hex"
	"time"
)

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
	GetID() string

	// GetIDShort returns the first 16 bytes of the task ID.
	GetIDShort() string
}

// Task defines the information of a task.
type Task struct {
	State   int           `json:"state,omitempty"`
	Plugin  string        `json:"plugin,omitempty"`
	ID      string        `json:"id,omitempty"`
	Node    string        `json:"node,omitempty"`
	Details interface{}   `json:"details,omitempty"`
	Log     []string      `json:"log,omitempty"`
	Time    time.Time     `json:"time,omitempty"`
	Buffer  *bytes.Buffer `json:"-"`
}

// GetID returns the ID of the task.
func (t *Task) GetID() string { return t.ID }

// GetIDShort returns the first 8 bytes of the task ID.
func (t *Task) GetIDShort() string { x, _ := hex.DecodeString(t.ID); return hex.EncodeToString(x[:8]) }
