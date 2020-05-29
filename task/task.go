package task

import (
	"fmt"
	"plugin"

	"github.com/chabad360/covey/task/types"
)

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

func getTask(identifier string) (types.ITask, bool) {
	if t, ok := tasks[identifier]; ok {
		return t, true
	} else if t, ok := tasksShort[identifier]; ok {
		return tasks[t], true
	}
	return nil, false
}
