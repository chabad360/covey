package main

import (
	"bytes"

	"github.com/chabad360/covey/task/types"
)

// Plugin is exposed to the module.
var Plugin plugin

// ShellTask fills the details field of the Task.
type ShellTask struct {
	Command    []string      `json:"command"`
	ExitStatus int           `json:"exit_status"`
	Buffer     *bytes.Buffer `json:"-"`
}

// Task implements the types.Task struct internally.
type Task struct {
	types.Task
	Details ShellTask `json:"details"`
}

type plugin struct{}

type newTaskInfo struct {
	Node    string   `json:"node"`
	Command []string `json:"command"`
	Plugin  string   `json:"plugin,omitempty"`
}
