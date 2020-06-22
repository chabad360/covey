package main

import (
	"fmt"

	"github.com/chabad360/covey/task/types"
	json "github.com/json-iterator/go"
)

// Plugin is exposed to the module.
var Plugin plugin

type plugin struct{}

// GetCommand returns the command to run on the node.
func (p *plugin) GetCommand(taskJSON []byte) (string, error) {
	var t types.Task
	if err := json.Unmarshal(taskJSON, &t); err != nil {
		return "", err
	}
	if t.Details["command"] == "" {
		return "", fmt.Errorf("shellPlugin: missing command")
	}

	return t.Details["command"], nil
}
