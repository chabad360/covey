package job

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/job/types"
	"github.com/go-playground/pure/v5"
	json "github.com/json-iterator/go"
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
		common.ErrorWriterCustom(w, fmt.Errorf("duplicate job: %v", job.Name), http.StatusBadRequest)
		return
	}

	if job.Cron != "" {
		if err := addCron(job.ID, job.Cron); err != nil {
			common.ErrorWriterCustom(w, err, http.StatusBadRequest)
		}
	}

	job.TaskHistory = []string{}
	job.ID = common.GenerateID(job)
	if err := AddJob(job); err != nil {
		common.ErrorWriter(w, err)
		return
	}

	common.Write(w, job)
}

func jobGet(w http.ResponseWriter, r *http.Request) {
	vars := pure.RequestVars(r)
	job, ok := GetJobWithTasks(vars.URLParam("job"))
	if !ok {
		common.ErrorWriter404(w, vars.URLParam("job"))
		return
	}
	common.Write(w, job)
}

func jobRun(w http.ResponseWriter, r *http.Request) {
	vars := pure.RequestVars(r)
	j, ok := GetJob(vars.URLParam("job"))
	if !ok {
		common.ErrorWriter404(w, vars.URLParam("job"))
		return
	}

	j.Run()
	if err := UpdateJob(*j); err != nil {
		log.Panic(err)
	}
}

// RegisterHandlers adds the handlers for the node module.
func RegisterHandlers(newRoute pure.IRouteGroup, singleRoute pure.IRouteGroup) {
	log.Println("Registering Job module API handlers...")
	newRoute.Post("/new", jobNew)

	j := singleRoute.Group("/:job")
	j.Post("", jobRun)
	j.Get("", jobGet)
}
