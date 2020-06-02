package types

import (
	"encoding/hex"
	"time"
)

// TaskState represents the current state of the task.
type TaskState int

// TaskPlugin defines the interface for Task module plugins.
type TaskPlugin interface {
	// NewTask creates a Task object and runs it in a go routine, and returns the Task.
	NewTask(taskJSON []byte) (ITask, error)

	// LoadTask loads the task.
	LoadTask(taskJSON []byte) (ITask, error)
}

// ITask Defines the read methods for a task.
type ITask interface {
	// GetLog returns the full log of the task.
	GetLog() []string

	// GetState returns the current state of the task.
	GetState() TaskState

	// GetID returns the ID of the task.
	GetID() string

	// GetIDShort returns the first 16 bytes of the task ID.
	GetIDShort() string
}

// Task defines the information of a task.
type Task struct {
	State   TaskState   `json:"state"`
	Plugin  string      `json:"plugin"`
	ID      string      `json:"id"`
	Node    string      `json:"node"`
	Details interface{} `json:"details"`
	Log     []string    `json:"log"`
	Time    time.Time   `json:"time"`
}

// GetID returns the ID of the task.
func (t *Task) GetID() string { return t.ID }

// GetIDShort returns the first 8 bytes of the task ID.
func (t *Task) GetIDShort() string { x, _ := hex.DecodeString(t.ID); return hex.EncodeToString(x[:8]) }

// GetState returns the current state of the task.
func (t *Task) GetState() TaskState { return t.State }
