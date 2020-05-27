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

var Plugin plugin

type newTaskInfo struct {
	Node    string   `json:"node,omitempty"`
	Command []string `json:"command,omitempty"`
	Plugin  string   `json:"plugin,omitempty"`
}

type ShellTask struct {
	Command    []string `json:"command,omitempty"`
	ExitStatus int      `json:"exit_status,omitempty"`
}

type Task struct {
	types.Task
	Details ShellTask `json:"details,omitempty"`
}

func (t *Task) GetUnread() []string {
	b := []byte{}
	copy(b, t.Buffer.Bytes())
	return proccessBytes(t, b)
}

func (t *Task) GetLog() []string {
	t.GetUnread()
	return t.Log
}

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
	n, err := node.GetNode(t.Node)
	if err != nil {
		return nil, err
	}

	b, c, err := n.Run(t.Details.Command)
	if err != nil {
		return nil, err
	}

	go func() {
		e := <-c
		if e == 0 {
			t.State = types.StateDone
		} else {
			t.State = types.StateError
			t.Details.ExitStatus = e
		}
	}()

	return b, nil
}

func proccessBytes(t *Task, b []byte) []string {
	c := []byte{}
	r := []string{}
	for _, bb := range b {
		if bb == '\n' {
			t.Log = append(t.Log, string(c))
			r = append(r, string(c))
			c = nil
		} else {
			c = append(c, bb)
		}
	}

	return r
}
