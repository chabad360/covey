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

	if err := AddJob(&job); err != nil {
		common.ErrorWriter(w, err)
	}

	common.Write(w, job)
}

func jobsGet(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()
	refreshDB()

	var q common.QueryParams
	err := q.Setup(r)
	common.ErrorWriter(w, err)

	var jobs interface{}

	if q.Expand {
		var j []models.Job
		err = q.Query("jobs", &j, db)
		jobs = j
	} else {
		var j []string
		err = q.Query("jobs", &j, db)
		jobs = j
	}
	common.ErrorWriter(w, err)

	common.Write(w, jobs)
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

	th, err := Run(j)
	common.ErrorWriter(w, err)

	common.Write(w, th)
}

func jobUpdate(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	vars := pure.RequestVars(r)
	job, ok := GetJob(vars.URLParam("job"))
	if !ok {
		common.ErrorWriter404(w, vars.URLParam("job"))
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(reqBody, &job); err != nil {
		common.ErrorWriterCustom(w, err, http.StatusBadRequest)
	}

	if job.Cron != "" {
		if err := addCron(job.ID, job.Cron); err != nil {
			common.ErrorWriterCustom(w, err, http.StatusBadRequest)
		}
	}

	err := UpdateJob(job)
	common.ErrorWriter(w, err)

	common.Write(w, job)
}

func jobDelete(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()
	refreshDB()

	vars := pure.RequestVars(r)
	j, ok := GetJob(vars.URLParam("job"))
	if !ok {
		common.ErrorWriter404(w, vars.URLParam("job"))
	}

	result := db.Delete(&j)
	if result.Error != nil {
		common.ErrorWriter(w, result.Error)
	}

	common.Write(w, vars.URLParam("job"))
}

// RegisterHandlers adds the handlers for the node module.
func RegisterHandlers(r pure.IRouteGroup) {
	log.Println("Registering Job module API handlers...")
	r.Post("", jobNew)
	r.Get("", jobsGet)

	j := r.Group("/:job")
	j.Post("", jobRun)
	j.Get("", jobGet)
	j.Put("", jobUpdate)
	j.Delete("", jobDelete)
}
