package task

import (
	"fmt"
	"github.com/chabad360/covey/models"
	"github.com/chabad360/covey/storage"
	json "github.com/json-iterator/go"
	"plugin"
)

// NewTask creates a new task.
func NewTask(taskJSON []byte) (*models.Task, error) {
	var t *models.Task
	if err := json.Unmarshal(taskJSON, &t); err != nil {
		return nil, err
	}

	_, ok := storage.GetNode(t.Node)
	if !ok {
		return nil, fmt.Errorf("%v is not a valid node", t.Node)
	}

	p, err := loadPlugin(t.Plugin)
	if err != nil {
		return nil, err
	}

	cmd, err := p.GetCommand(taskJSON)
	if err != nil {
		return nil, err
	}

	t.Command = cmd

	err = storage.AddTask(t)
	if err != nil {
		return nil, err
	}

	err = queueTask(t.Node, t.ID, t.Command)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func loadPlugin(pluginName string) (TaskPlugin, error) {
	p, err := plugin.Open("./plugins/task/" + pluginName + ".so")
	if err != nil {
		return nil, err
	}

	n, err := p.Lookup("Plugin")
	if err != nil {
		return nil, err
	}

	var s TaskPlugin

	s, ok := n.(TaskPlugin)
	if !ok {
		return nil, fmt.Errorf(pluginName, " does not provide a TaskPlugin")
	}

	return s, nil
}
