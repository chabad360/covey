package job

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/chabad360/covey/common"
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
