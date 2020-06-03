package types

import (
	"encoding/hex"
	"encoding/json"
	"log"

	"github.com/chabad360/covey/task"
)

type Job struct {
	Name        string             `json:"name"`
	ID          string             `json:"id"`
	Cron        string             `json:"cron,omitempty"`
	Nodes       []string           `json:"nodes"`
	Tasks       map[string]jobTask `json:"tasks"`
	TaskHistory []string           `json:"task_history,omitempty"`
}

type JobWithTasks struct {
	Job
	TaskHistory []interface{} `json:"task_history,omitempty"`
}

type jobTask struct {
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
	for z := range j.Tasks {
		t := j.Tasks[z]
		for node := range j.Nodes {
			t.Node = j.Nodes[node]
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
