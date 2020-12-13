package shell

import (
	"encoding/json"
	"fmt"
	"github.com/chabad360/covey/models/safe"
)

// Plugin is exposed to the module
type Plugin struct{}

// GetCommand returns the command to run on the node.
func (p *Plugin) GetCommand(task safe.Task) (string, error) {
	if task.Details["command"] == "" {
		return "", fmt.Errorf("shellPlugin: missing command")
	}
	return task.Details["command"], nil
}

// GetInputs returns the input for the Shell plugin.
func (p *Plugin) GetInputs() safe.Form {
	return safe.Form{
		Inputs: []safe.Input{
			{
				Name:     "command",
				Label:    "Command",
				Type:     safe.Text,
				Required: true,
			},
		},
	}
}

// GetFetchCommand returns the command and callback to get basic info about the node.
func (p *Plugin) GetFetchCommand() (string, func([]string) ([]byte, error)) {
	return "uname -s && uname -n && uname -r && uname -m && uname -o", func(output []string) ([]byte, error) {
		return json.Marshal(struct {
			KernelName      string `json:"kernel-name"`
			HostName        string `json:"hostname"`
			KernelRelease   string `json:"kernel-release"`
			Machine         string `json:"machine"`
			OperatingSystem string `json:"operatingSystem"`
		}{
			KernelName:      output[0],
			HostName:        output[1],
			KernelRelease:   output[2],
			Machine:         output[3],
			OperatingSystem: output[4],
		})
	}
}
