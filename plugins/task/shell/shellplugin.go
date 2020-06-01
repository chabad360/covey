package main

import (
	"encoding/json"
	"fmt"
	"log"
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

	b, err := runTask(&t)
	if err != nil {
		return nil, err
	}
	t.Details.Buffer = b

	id, err := common.GenerateID(t)
	if err != nil {
		return nil, err
	}
	t.ID = *id

	return &t, nil
}

func (p *plugin) LoadTask(taskJSON []byte) (types.ITask, error) {
	var t Task
	if err := json.Unmarshal(taskJSON, &t); err != nil {
		return nil, err
	}
	log.Println("Loading task", t.ID)

	return &t, nil
}
