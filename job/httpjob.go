package job

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/storage"
	"github.com/gorilla/mux"
)

func jobNew(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var job Job
	reqBody, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(reqBody, &job); err != nil {
		common.ErrorWriter(w, err)
		return
	}

	if _, ok := GetJob(job.Name); ok {
		common.ErrorWriter(w, fmt.Errorf("Duplicate job: %v", job.Name))
		return
	}

	if job.Cron != "" {
		_, err := cronTab.AddFunc(job.Cron, func() func() {
			return func() { // This little bundle of joy allows the job to occur despite not being an object.
				j, _ := GetJob(job.GetID())
				j.Run()
			}
		}())
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
	// jobs[job.GetID()] = &job
	// jobsShort[job.GetIDShort()] = job.GetID()
	// jobsName[job.GetName()] = job.GetID()

	if err = storage.AddItem("jobs", job.GetID(), job.GetIDShort(), job); err != nil {
		common.ErrorWriter(w, err)
		return
	}

	z, err := storage.GetItem("jobs", job.GetID(), job)
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}

	j, err := json.MarshalIndent(z, "", "\t")
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
	job, ok := GetJobWithTasks(vars["job"])
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 %v not found", vars["job"])
		return
	}

	// for _, task := range job.TaskHistory {
	// 	task.GetLog()
	// }

	// jobs[job.GetID()] = job

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
	j, ok := GetJob(vars["job"])
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 %v not found", vars["job"])
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
