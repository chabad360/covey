package main

import (
	"bytes"

	"github.com/chabad360/covey/task/types"
)

// Plugin is exposed to the module.
var Plugin plugin

// Task implements the types.Task struct internally.
type Task struct {
	types.Task
	Buffer *bytes.Buffer `json:"-"`
}

type plugin struct{}

type newTaskInfo struct {
	Node    string   `json:"node"`
	Command []string `json:"command"`
	Plugin  string   `json:"plugin,omitempty"`
}
