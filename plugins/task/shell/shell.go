package main

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/chabad360/covey/node"
	"github.com/chabad360/covey/task"
	"github.com/chabad360/covey/task/types"
)

func runTask(t *Task) (*bytes.Buffer, error) {
	n, ok := node.GetNode(t.Node)
	if !ok {
		return nil, fmt.Errorf("%v is not a valid node", t.Node)
	}

	b, c, err := n.Run([]string{t.Details["command"]})
	if err != nil {
		return nil, err
	}
	t.State = types.StateRunning

	go func() {
		e := <-c
		if e == 0 {
			t.State = types.StateDone
		} else {
			t.State = types.StateError
		}
		t.Details["exit_status"] = strconv.Itoa(e)
		t.GetLog()
		task.SaveTask(t)
	}()

	return b, nil
}
