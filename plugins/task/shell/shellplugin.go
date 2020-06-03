package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/task/types"
)

// NewTask creates a new task.
func (p *plugin) NewTask(taskJSON []byte) (types.ITask, error) {
	var t Task
	if err := json.Unmarshal(taskJSON, &t); err != nil {
		return nil, err
	}
	if t.Details.Command == nil {
		return nil, fmt.Errorf("Missing command")
	}

	t.Details.ExitStatus = 256
	t.Log = []string{}
	t.State = types.StateStarting
	t.Time = time.Now()
	t.ID = common.GenerateID(t)

	b, err := runTask(&t)
	if err != nil {
		return nil, err
	}
	t.Details.Buffer = b

	return &t, nil
}
