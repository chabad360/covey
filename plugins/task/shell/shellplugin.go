package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/task/types"
	json "github.com/json-iterator/go"
)

// NewTask creates a new task.
func (p *plugin) NewTask(taskJSON []byte) (types.ITask, error) {
	var t Task
	if err := json.Unmarshal(taskJSON, &t); err != nil {
		return nil, err
	}
	if t.Details["command"] == "" {
		return nil, fmt.Errorf("shellPlugin: missing command")
	}

	t.Details["exit_status"] = strconv.Itoa(256)
	t.Log = []string{}
	t.State = types.StateStarting
	t.Time = time.Now()
	t.ID = common.GenerateID(t)

	b, err := runTask(&t)
	if err != nil {
		return nil, err
	}
	t.Buffer = b

	return &t, nil
}
