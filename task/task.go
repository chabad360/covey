package task

import (
	"fmt"
	"log"
	"plugin"
	"time"

	json "github.com/json-iterator/go"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/node"
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

	t.ExitCode = 256
	t.Log = []string{}
	t.State = types.StateQueued
	t.Time = time.Now()
	t.Command = cmd
	t.ID = common.GenerateID(t)

	tasks[t.GetID()] = t
	tasksShort[t.GetIDShort()] = t.GetID()
	SaveTask(t)

	go runTask(t)
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
	// If the task is still running, return it instead of the db version.
	if x, ok := tasks[identifier]; ok {
		return x, true
	} else if x, ok := tasksShort[identifier]; ok {
		return tasks[x], true
	}

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
func SaveTask(t *types.Task) {
	// Only useful if the task is still running, i.e. it's in the tasks map.
	if _, ok := tasks[t.GetID()]; !ok {
		return
	}
	if _, err := getTaskJSON(t.GetID()); err != nil { // If the task isn't in the database yet
		if err = addTask(t); err != nil {
			log.Println(err)
		}
	} else { // Otherwise:
		if err = updateTask(t); err != nil {
			log.Println(err)
		}
	}

	// Update the task in the tasks map or remove it if it's done.
	if state := t.GetState(); !(state == types.StateRunning || state == types.StateQueued) {
		delete(tasks, t.GetID())
		delete(tasksShort, t.GetIDShort())
	} else {
		tasks[t.GetID()] = t
	}
}

func runTask(t *types.Task) error {
	n, ok := node.GetNode(t.Node)
	if !ok {
		return fmt.Errorf("%v is not a valid node", t.Node)
	}

	b, c, err := n.Run([]string{t.Command})
	if err != nil {
		return err
	}
	t.State = types.StateRunning
	t.Buffer = b

	e := <-c
	if e == 0 {
		t.State = types.StateDone
	} else {
		t.State = types.StateError
	}
	t.ExitCode = e
	t.GetLog()
	SaveTask(t)

	return nil
}
