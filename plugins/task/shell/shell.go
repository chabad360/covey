package main

import (
	"fmt"
	"github.com/chabad360/covey/models"
)

// Plugin is exposed to the module.
var Plugin plugin

type plugin struct{}

// GetCommand returns the command to run on the node.
func (p *plugin) GetCommand(task models.Task) (string, error) {
	if task.Details["command"] == "" {
		return "", fmt.Errorf("shellPlugin: missing command")
	}
	return task.Details["command"], nil
}
