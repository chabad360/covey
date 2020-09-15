package task

import (
	"container/list"
	"github.com/chabad360/covey/models"
	json "github.com/json-iterator/go"
)

// AgentTask contains the information about a task that is sent to an agent.
type agentTask struct {
	Command string `json:"command"`
	ID      string `json:"id"`
}

// List implements List type as well as the json.Marshaler interface.
type List struct{ list.List }

// MarshalJSON implements the json.Marshaler interface.
func (l *List) MarshalJSON() ([]byte, error) {
	m := make(map[int]agentTask)
	i := 0
	for e := l.Front(); e != nil; e = e.Next() {
		m[i] = e.Value.(agentTask)
		i++
	}
	return json.Marshal(m)
}

// Plugin defines the interface for Task module plugins.
type Plugin interface {
	// GetCommand returns the command to run the server.
	GetCommand(task models.Task) (string, error)

	// GetFetchCommand returns a command to run which will be used to fetch relevant information about the node, and a callback that returns JSON metadata to send the output too.
	GetFetchCommand() (string, func([]string) ([]byte, error))

	// GetInputs takes the JSON metadata and returns a JSON that can be converted to a ui.Form object.
	GetInputs([]byte) ([]byte, error)

	// GetUUID returns the UUID of the plugin.
	GetUUID() string
}
