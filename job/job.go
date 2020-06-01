package job

import (
	"fmt"
	"log"
	"os"

	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/task"
	"github.com/chabad360/covey/task/types"
	"github.com/gorilla/mux"
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

// type IJob interface {
// 	GetName() string

// 	GetID() string

// 	GetIDShort() string

// 	Run()
// }

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

func jobNew(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var job Job
	reqBody, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(reqBody, &job); err != nil {
		common.ErrorWriter(w, err)
		return
	}

	if _, ok := getJob(job.Name); ok {
		common.ErrorWriter(w, fmt.Errorf("Duplicate job: %v", job.Name))
		return
	}

	if job.Cron != "" {
		_, err := cronTab.AddFunc(job.Cron, job.Run)
		if err != nil {
			common.ErrorWriter(w, err)
			return
		}
	}

	id, err := common.GenerateID(job)
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}

	job.ID = *id
	jobs[job.GetID()] = &job
	jobsShort[job.GetIDShort()] = job.GetID()
	jobsName[job.GetName()] = job.GetID()
	j, err := json.MarshalIndent(jobs, "", "  ")
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}
	f, err := os.Create("./config/jobs.json")
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}
	defer f.Close()
	if err = f.Chmod(0600); err != nil {
		common.ErrorWriter(w, err)
		return
	}
	if _, err = f.Write(j); err != nil {
		common.ErrorWriter(w, err)
		return
	}

	j, err = json.MarshalIndent(job, "", "  ")
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}

	w.Header().Set("Location", "/api/v1/job/"+job.GetIDShort())
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(j))
}

func jobGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	job, ok := getJob(vars["job"])
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 %v not found", vars["job"])
		return
	}

	for _, task := range job.TaskHistory {
		task.GetLog()
	}

	jobs[job.GetID()] = job

	w.Header().Add("Content-Type", "application/json")
	j, err := json.MarshalIndent(job, "", "\t")
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(j))
}

func jobRun(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	j, ok := getJob(vars["job"])
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "404 %v not found", vars["job"])
		return
	}

	j.Run()
}

// RegisterHandlers adds the mux handlers for the node module.
func RegisterHandlers(r *mux.Router) {
	log.Println("Registering Job module API handlers...")
	r.HandleFunc("/new", jobNew).Methods("POST")
	r.HandleFunc("/{job}", jobGet).Methods("GET")
	r.HandleFunc("/{job}", jobRun).Methods("POST")

}

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
