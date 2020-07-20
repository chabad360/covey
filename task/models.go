package task

import (
	"container/list"
	json "github.com/json-iterator/go"
)

// AgentTask contains the information about a task that is send to an agent.
type agentTask struct {
	Command string `json:"command"`
	ID      string `json:"id"`
}

// TaskList implements List type as well as the json.Marshaler interface.
type TaskList struct{ list.List }

// TaskInfo contains new information about a running task.
type TaskInfo struct {
	Log      []string `json:"log"`
	ExitCode int      `json:"exit_code"`
	ID       string   `json:"id"`
}

// MarshalJSON implements the json.Marshaler interface.
func (l *TaskList) MarshalJSON() ([]byte, error) {
	m := make(map[int]agentTask)
	i := 0
	for e := l.Front(); e != nil; e = e.Next() {
		m[i] = e.Value.(agentTask)
		i++
	}
	return json.Marshal(m)
}

// TaskPlugin defines the interface for Task module plugins.
type TaskPlugin interface {
	// GetCommand returns the command to run the server.
	GetCommand(taskJSON []byte) (string, error)
}
