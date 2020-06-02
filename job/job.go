package job

import (
	"log"

	"encoding/hex"
	"encoding/json"

	"github.com/chabad360/covey/task"
	"github.com/chabad360/covey/task/types"
	"github.com/robfig/cron/v3"
)

var (
	jobs      = make(map[string]*Job)
	jobsShort = make(map[string]string)
	jobsName  = make(map[string]string)
	cronTab   = cron.New()
)

type Job struct {
	Name        string                 `json:"name"`
	ID          string                 `json:"id"`
	Cron        string                 `json:"cron,omitempty"`
	Nodes       []string               `json:"nodes"`
	Tasks       map[string]jobTask     `json:"tasks"`
	TaskHistory map[string]types.ITask `json:"task_history,omitempty"`
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
			if j.TaskHistory == nil {
				j.TaskHistory = make(map[string]types.ITask)
			}
			j.TaskHistory[r.GetID()] = r
		}
	}
}

// LoadConfig loads up the configuration and starts the cronTab.
func LoadConfig() {
	cronTab.Start()
}

// getJob checks if a job with the identifier exists and returns it.
func getJob(identifier string) (*Job, bool) {
	if j, ok := jobs[identifier]; ok {
		return j, true
	} else if j, ok := jobsShort[identifier]; ok {
		t := jobs[j]
		return t, true
	} else if j, ok := jobsName[identifier]; ok {
		t := jobs[j]
		return t, true
	}
	return nil, false
}
