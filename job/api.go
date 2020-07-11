package job

import (
	"fmt"
	"github.com/chabad360/covey/models"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/chabad360/covey/common"
	"github.com/go-playground/pure/v5"
	json "github.com/json-iterator/go"
)

func jobNew(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	w.Header().Set("Content-Type", "application/json")

	var job models.Job

	reqBody, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(reqBody, &job); err != nil {
		common.ErrorWriterCustom(w, err, http.StatusBadRequest)
	}

	if _, ok := GetJob(job.Name); ok {
		common.ErrorWriterCustom(w, fmt.Errorf("duplicate job: %v", job.Name), http.StatusBadRequest)
	}

	if job.Cron != "" {
		if err := addCron(job.ID, job.Cron); err != nil {
			common.ErrorWriterCustom(w, err, http.StatusBadRequest)
		}
	}

	if err := AddJob(job); err != nil {
		common.ErrorWriter(w, err)
	}

	common.Write(w, job)
}

func jobGet(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	vars := pure.RequestVars(r)
	job, ok := GetJob(vars.URLParam("job"))
	if !ok {
		common.ErrorWriter404(w, vars.URLParam("job"))
	}

	common.Write(w, job)
}

func jobRun(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	vars := pure.RequestVars(r)
	j, ok := GetJob(vars.URLParam("job"))
	if !ok {
		common.ErrorWriter404(w, vars.URLParam("job"))
	}

	Run(j)

	if err := UpdateJob(*j); err != nil {
		log.Panic(err)
	}
}

// RegisterHandlers adds the handlers for the node module.
func RegisterHandlers(r pure.IRouteGroup) {
	log.Println("Registering Job module API handlers...")
	r.Post("", jobNew)

	j := r.Group("/:job")
	j.Post("", jobRun)
	j.Get("", jobGet)
}
