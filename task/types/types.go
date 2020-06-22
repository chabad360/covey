package types

import (
	"bytes"
	"encoding/hex"
	"time"
)

// TaskState represents the current state of the task.
type TaskState int

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
	Buffer   *bytes.Buffer     `json:"-"`
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
func (t *Task) GetLog() []string {
	if t.Buffer != nil { // Ensure buffer exists
		b := t.Buffer.Bytes()

		var line []byte
		var log []string

		for _, bb := range b { // For each byte...
			if bb == '\n' { // If that byte is a newline:
				log = append(log, string(line)) // Add that line to the log
				line = nil                      // And start the next one
			} else { // Otherwise,
				line = append(line, bb) // Add it to the line
			}
		}

		if len(line) > 0 { // If the last line didn't end with a newline
			log = append(log, string(line)) // Append it
		}

		if len(log) > 0 { // Only set the log if there is stuff on it, otherwise we get empty logs.
			t.Log = log
		}

		t.Buffer.Reset() // Finally, reset the buffer.
	}

	return t.Log
}
