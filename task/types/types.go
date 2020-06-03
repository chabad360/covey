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

// ITask Defines the read methods for a task.
type ITask interface {
	// GetLog returns the full log of the task.
	GetLog() []string

	// GetState returns the current state of the task.
	GetState() TaskState

	// GetPlugin returns the plugin of the task.
	GetPlugin() string

	// GetNode returns the node of the task..
	GetNode() string

	// GetTime returns the time of the task.
	GetTime() time.Time

	// GetDetails returns the details of the task.
	GetDetails() interface{}

	// GetID returns the ID of the task.
	GetID() string

	// GetIDShort returns the first 16 bytes of the task ID.
	GetIDShort() string
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
func (t *Task) GetDetails() interface{} { return t.Details }

// GetLog returns the log of the task.
func (t *Task) GetLog() []string { return t.Log }
