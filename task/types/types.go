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

	// GetID returns the ID of the task.
	GetID() string

	// GetIDShort returns the first 16 bytes of the task ID.
	GetIDShort() string
}

// Task defines the information of a task.
type Task struct {
	State   int           `json:"state"`
	Plugin  string        `json:"plugin"`
	ID      string        `json:"id"`
	Node    string        `json:"node"`
	Details interface{}   `json:"details"`
	Log     []string      `json:"log"`
	Time    time.Time     `json:"time"`
	Buffer  *bytes.Buffer `json:"-"`
}

// GetID returns the ID of the task.
func (t *Task) GetID() string { return t.ID }

// GetIDShort returns the first 8 bytes of the task ID.
func (t *Task) GetIDShort() string { x, _ := hex.DecodeString(t.ID); return hex.EncodeToString(x[:8]) }
