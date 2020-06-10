package types

import (
	"encoding/hex"
	"encoding/json"
	"log"

	"github.com/chabad360/covey/task"
)

// Job contains the information for a given job
type Job struct {
	Name        string             `json:"name"`
	ID          string             `json:"id"`
	Cron        string             `json:"cron,omitempty"`
	Nodes       []string           `json:"nodes"`
	Tasks       map[string]JobTask `json:"tasks"`
	TaskHistory []string           `json:"task_history"`
}

// JobWithTasks is the same thing as a Job but with the task history enumerated.
type JobWithTasks struct {
	Job
	TaskHistory []interface{} `json:"task_history,omitempty"`
}

type JobTask struct {
	Plugin  string      `json:"plugin"`
	Details interface{} `json:"details"`
	Node    string      `json:"node,omitempty"`
}

// GetName returns the name of the job.
func (j *Job) GetName() string { return j.Name }

// GetID returns the ID of the job.
func (j *Job) GetID() string { return j.ID }

// GetIDShort returns the first 8 bytes of the job ID.
func (j *Job) GetIDShort() string { x, _ := hex.DecodeString(j.ID); return hex.EncodeToString(x[:8]) }

// Run runs each task in succession on the specified nodes (concurrently).
func (j *Job) Run() {
	for _, t := range j.Tasks {
		for _, node := range j.Nodes {
			t.Node = node
			x, err := json.Marshal(t)
			if err != nil {
				log.Panic(err)
			}

			r, err := task.NewTask(x)
			if err != nil {
				log.Panic(err)
			}

			j.TaskHistory = append(j.TaskHistory, r.GetID())
		}
	}
}
