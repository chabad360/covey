package task

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"plugin"

	"github.com/chabad360/covey/storage"
	"github.com/chabad360/covey/task/types"
)

// NewTask creates a new task.
func NewTask(taskJSON []byte) (types.ITask, error) {
	var task types.Task
	if err := json.Unmarshal(taskJSON, &task); err != nil {
		return nil, err
	}

	p, err := loadPlugin(task.Plugin)
	if err != nil {
		return nil, err
	}

	t, err := p.NewTask(taskJSON)
	if err != nil {
		return nil, err
	}

	tasks[t.GetID()] = t
	tasksShort[t.GetIDShort()] = t.GetID()
	SaveTask(t)
	return t, nil
}

// LoadConfig loads up the stored nodes
func LoadConfig() {
	log.Println("Loading Task Config")
	f, err := os.Open("./config/tasks.json")
	if err != nil {
		log.Println("Error loading task config")
		return
	}
	defer f.Close()
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
func GetTask(identifier string) (types.ITask, bool) {
	var t *types.Task
	n, err := storage.GetItem("tasks", identifier, t)
	if err != nil {
		log.Println(err)
		return nil, false
	}
	j, err := json.Marshal(n)
	if err != nil {
		log.Println(err)
		return nil, false
	}
	x, err := loadTask(j)
	if err != nil {
		log.Println(err)
		return nil, false
	}

	// If the task is still running, return it instead of the db version.
	if x, ok := tasks[identifier]; ok {
		return x, true
	} else if x, ok := tasksShort[identifier]; ok {
		return tasks[x], true
	}

	return x, true
}

func loadTask(taskJSON []byte) (types.ITask, error) {
	var z types.Task
	if err := json.Unmarshal(taskJSON, &z); err != nil {
		return nil, err
	}
	p, err := loadPlugin(z.Plugin)
	if err != nil {
		return nil, err
	}
	t, err := p.LoadTask(taskJSON)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// SaveTask saves a live task to the database.
func SaveTask(t types.ITask) {
	// Only useful if the task is still running, i.e. it's in the tasks map.
	if _, ok := tasks[t.GetID()]; !ok {
		return
	}

	var z *types.Task
	_, err := storage.GetItem("tasks", t.GetID(), z)
	if err != nil { // If the task isn't in the database yet
		log.Println(err)
		if err = storage.AddItem("tasks", t.GetID(), t.GetIDShort(), t); err != nil {
			log.Println(err)
		}
	} else { // Otherwise:
		if err = storage.UpdateItem("tasks", t.GetID(), t); err != nil {
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
