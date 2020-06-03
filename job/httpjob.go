package job

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/job/types"
	"github.com/gorilla/mux"
)

func jobNew(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var job types.Job
	reqBody, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(reqBody, &job); err != nil {
		common.ErrorWriterCustom(w, err, http.StatusBadRequest)
		return
	}

	if _, ok := GetJob(job.Name); ok {
		common.ErrorWriterCustom(w, fmt.Errorf("Duplicate job: %v", job.Name), http.StatusBadRequest)
		return
	}

	if job.Cron != "" {
		if err := addCron(job.ID, job.Cron); err != nil {
			common.ErrorWriterCustom(w, err, http.StatusBadRequest)
		}
	}

	job.ID = common.GenerateID(job)
	// jobs[job.GetID()] = &job
	// jobsShort[job.GetIDShort()] = job.GetID()
	// jobsName[job.GetName()] = job.GetID()

	if err := AddJob(job); err != nil {
		common.ErrorWriter(w, err)
		return
	}

	common.Write(w, job)
}

func jobGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	job, ok := GetJobWithTasks(vars["job"])
	if !ok {
		common.ErrorWriter404(w, vars["job"])
		return
	}
	common.Write(w, job)
}

func jobRun(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	j, ok := GetJob(vars["job"])
	if !ok {
		common.ErrorWriter404(w, vars["job"])
		return
	}

	j.Run()
	if err := UpdateJob(*j); err != nil {
		log.Panic(err)
	}
}

// RegisterHandlers adds the mux handlers for the node module.
func RegisterHandlers(r *mux.Router) {
	log.Println("Registering types.Job module API handlers...")
	r.HandleFunc("/new", jobNew).Methods("POST")
	r.HandleFunc("/{job}", jobGet).Methods("GET")
	r.HandleFunc("/{job}", jobRun).Methods("POST")

}
