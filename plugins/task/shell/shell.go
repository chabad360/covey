package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/node"
	"github.com/chabad360/covey/task/types"
)

type plugin struct{}

// Plugin is exposed to the module.
var Plugin plugin

type newTaskInfo struct {
	Node    string   `json:"node"`
	Command []string `json:"command"`
	Plugin  string   `json:"plugin,omitempty"`
}

type ShellTask struct {
	Command    []string `json:"command"`
	ExitStatus int      `json:"exit_status"`
}

type Task struct {
	types.Task
	Details ShellTask `json:"details"`
}

// GetLog reads the unread buffer and adds it to the task's log, then returns that log.
func (t *Task) GetLog() []string {
	// There is a reason why some tty clients print every rewrite of a line,
	// it's much easier to simply process every line as a new line
	// but if we can store the fact that \n hasn't yet been given,
	// we can solve that issue.
	// Also escaping is another issue...
	b := t.Buffer.Bytes()
	c := []byte{}
	l := []string{}
	for _, bb := range b {
		if bb == '\n' {
			l = append(l, string(c))
			c = nil
		} else {
			c = append(c, bb)
		}
	}
	if len(c) > 0 {
		l = append(l, string(c))
	}
	if len(l) > 0 {
		t.Log = l
	}
	t.Buffer.Reset()
	return t.Log
}

// NewTask creates a new task.
func (p *plugin) NewTask(taskJSON []byte) (types.ITask, error) {
	var taskInfo newTaskInfo
	if err := json.Unmarshal(taskJSON, &taskInfo); err != nil {
		return nil, err
	}
	if taskInfo.Command == nil {
		return nil, fmt.Errorf("Missing command")
	}

	t := Task{
		Details: ShellTask{
			Command: taskInfo.Command,
		},
	}
	t.Node = taskInfo.Node
	t.Log = []string{}
	t.Plugin = taskInfo.Plugin
	t.State = types.StateStarting
	t.Time = time.Now()
	t.Details.ExitStatus = 256

	b, err := runTask(&t)
	if err != nil {
		return nil, err
	}
	t.Buffer = b

	id, err := common.GenerateID(t)
	if err != nil {
		return nil, err
	}
	t.ID = *id

	return &t, nil
}

func runTask(t *Task) (*bytes.Buffer, error) {
	n, ok := node.GetNode(t.Node)
	if !ok {
		return nil, fmt.Errorf("%v is not a valid node", t.Node)
	}

	b, c, err := n.Run(t.Details.Command)
	if err != nil {
		return nil, err
	}

	go func() {
		e := <-c
		if e == 0 {
			t.State = types.StateDone
			t.Details.ExitStatus = e
		} else {
			t.State = types.StateError
			t.Details.ExitStatus = e
		}
	}()

	return b, nil
}
