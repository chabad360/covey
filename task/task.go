package task

import (
	"fmt"
	"log"
	"plugin"

	json "github.com/json-iterator/go"

	"github.com/chabad360/covey/task/types"
)

// NewTask creates a new task.
func NewTask(taskJSON []byte) (types.ITask, error) {
	var t types.Task
	if err := json.Unmarshal(taskJSON, &t); err != nil {
		return nil, err
	}

	p, err := loadPlugin(t.Plugin)
	if err != nil {
		return nil, err
	}

	task, err := p.NewTask(taskJSON)
	if err != nil {
		return nil, err
	}

	tasks[task.GetID()] = task
	tasksShort[task.GetIDShort()] = task.GetID()
	SaveTask(task)
	return task, nil
}

// LoadConfig loads up the stored nodes
// func LoadConfig() {
// 	log.Println("Placeholder")
// }

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
func GetTask(identifier string) (types.ITask, bool) {
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
func SaveTask(t types.ITask) {
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
	if state := t.GetState(); !(state == types.StateRunning || state == types.StateStarting) {
		delete(tasks, t.GetID())
		delete(tasksShort, t.GetIDShort())
	} else {
		tasks[t.GetID()] = t
	}
}
