package task

import (
	"fmt"
	"log"
	"plugin"
	"time"

	json "github.com/json-iterator/go"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/task/types"
)

// NewTask creates a new task.
func NewTask(taskJSON []byte) (*types.Task, error) {
	var t *types.Task
	if err := json.Unmarshal(taskJSON, &t); err != nil {
		return nil, err
	}

	p, err := loadPlugin(t.Plugin)
	if err != nil {
		return nil, err
	}

	cmd, err := p.GetCommand(taskJSON)
	if err != nil {
		return nil, err
	}

	t.ExitCode = 258
	t.Log = []string{}
	t.State = types.StateQueued
	t.Time = time.Now()
	t.Command = cmd
	t.ID = common.GenerateID(t)
	addTask(t)

	err = QueueTask(t.Node, t.ID, t.Command)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func loadPlugin(pluginName string) (types.TaskPlugin, error) {
	p, err := plugin.Open("./plugins/task/" + pluginName + ".so")
	if err != nil {
		return nil, err
	}

	n, err := p.Lookup("Plugin")
	if err != nil {
		return nil, err
	}

	var s types.TaskPlugin
	s, ok := n.(types.TaskPlugin)
	if !ok {
		return nil, fmt.Errorf(pluginName, " does not provide a TaskPlugin")
	}

	return s, nil
}

// GetTask checks if a task with the identifier exists and returns it.
func GetTask(identifier string) (*types.Task, bool) {
	t, err := getTaskJSON(identifier)
	if err != nil {
		return nil, false
	}
	var x types.Task
	if err := json.Unmarshal(t, &x); err != nil {
		return nil, false
	}
	return &x, true
}

// SaveTask saves a live task to the database.
func SaveTask(t *types.TaskInfo) {
	task, ok := GetTask(t.ID)
	if !ok {
		return
	}
	if task.ExitCode != t.ExitCode && t.Log != nil {
		task.ExitCode = t.ExitCode
		if t.ExitCode == 0 {
			task.State = types.StateDone
		} else if t.ExitCode == 257 {
			task.State = types.StateRunning
		} else {
			task.State = types.StateError
		}
		task.Log = append(task.Log, t.Log...)

		if err := updateTask(task); err != nil {
			log.Println(err)
		}
	}
}
