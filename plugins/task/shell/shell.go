package main

import (
	"bytes"
	"fmt"

	"github.com/chabad360/covey/node"
	"github.com/chabad360/covey/task/types"
)

func runTask(t *Task) (*bytes.Buffer, error) {
	n, ok := node.GetNode(t.Node)
	if !ok {
		return nil, fmt.Errorf("%v is not a valid node", t.Node)
	}

	b, c, err := n.Run(t.Details.Command)
	if err != nil {
		return nil, err
	}
	t.State = types.StateRunning

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
