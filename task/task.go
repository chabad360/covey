package task

import (
	"fmt"

	json "github.com/json-iterator/go"

	"github.com/chabad360/covey/models"
	"github.com/chabad360/covey/plugins"
	"github.com/chabad360/covey/storage"
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

	p, err := plugins.GetTaskPlugin(t.Plugin)
	if err != nil {
		return nil, err
	}

	cmd, err := p.GetCommand(t.ToSafe())
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
