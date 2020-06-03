package job

import (
	"log"

	"encoding/hex"
	"encoding/json"

	"github.com/chabad360/covey/storage"
	"github.com/chabad360/covey/task"
	"github.com/robfig/cron/v3"
)

var (
	// jobs      = make(map[string]*Job)
	// jobsShort = make(map[string]string)
	// jobsName  = make(map[string]string)
	cronTab = cron.New()
)

// IJob defines the interface for a job.
// Would'a been nice not to need this...
// type IJob interface {
// 	// GetName returns the name of the job.
// 	GetName() string

// 	// GetID returns the ID of the job.
// 	GetID() string

// 	// GetIDShort returns the first 8 bytes of the job ID.
// 	GetIDShort() string

// 	// Run executes the the job.
// 	Run()
// }

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
	if err := storage.UpdateItem("jobs", j.GetID(), j); err != nil {
		log.Panic(err)
	}
}

// LoadConfig loads up the configuration and starts the cronTab.
func LoadConfig() {
	cronTab.Start()
}

// GetJob checks if a job with the identifier exists and returns it.
func GetJob(identifier string) (*Job, bool) {
	var t Job
	_, err := storage.GetItem("jobs", identifier)
	if err != nil {
		log.Println(err)
		return nil, false
	}

	return &t, true
}

// GetJobWithTasks checks if a job with the identifier exists and returns it along with its tasks.
func GetJobWithTasks(identifier string) (*JobWithTasks, bool) {
	var t JobWithTasks
	_, err := storage.GetJob(identifier, &t)
	if err != nil {
		log.Println(err)
		return nil, false
	}

	return &t, true
}
