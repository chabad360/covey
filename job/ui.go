package job

import (
	"fmt"
	"github.com/chabad360/covey/models"
	"net/http"
	"strings"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/ui"
	"github.com/go-playground/pure/v5"
)

func uiJobs(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()
	refreshDB()

	var jobs []models.Job
	result := db.Find(&jobs)
	common.ErrorWriter(w, result.Error)

	p := &ui.Page{
		Title:   "Jobs",
		URL:     strings.Split(r.URL.Path, "/"),
		Details: struct{ Jobs []models.Job }{Jobs: jobs},
	}

	t := ui.GetTemplate("jobsAll")
	err := t.ExecuteTemplate(w, "base", p)
	common.ErrorWriter(w, err)
}

func uiJobSingle(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	vars := pure.RequestVars(r)

	job, ok := getJobWithFullHistory(vars.URLParam("job"))
	common.ErrorWriter404(w, vars.URLParam("job"), ok)

	if r.URL.Query().Get("run") == "true" {
		j, _ := getJob(vars.URLParam("job"))
		_, err := run(j)
		common.ErrorWriter(w, err)

		job, _ = getJobWithFullHistory(vars.URLParam("job"))
	}

	p := &ui.Page{
		Title:   fmt.Sprintf("Job %s", vars.URLParam("job")),
		URL:     strings.Split(r.URL.Path, "/"),
		Details: struct{ Job *models.JobWithTasks }{job},
	}

	t := ui.GetTemplate("jobsSingle")
	err := t.ExecuteTemplate(w, "base", p)
	common.ErrorWriter(w, err)
}

// UIJobNew returns the form for creating a new task.
func UIJobNew(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()
	refreshDB()

	var nodes []string
	result := db.Table("nodes").Select("name").Scan(&nodes)
	common.ErrorWriter(w, result.Error)

	p := &ui.Page{
		Title: "New Job",
		URL:   strings.Split(r.URL.Path, "/"),
		Details: struct {
			Plugins []string
			Nodes   []string
		}{[]string{"shell"}, nodes},
	}

	t := ui.GetTemplate("jobsNew")
	err := t.ExecuteTemplate(w, "base", p)
	common.ErrorWriter(w, err)
}

// RegisterUIHandlers registers the HTTP handlers for the jobs UI.
func RegisterUIHandlers(r pure.IRouteGroup) {
	r.Get("", uiJobs)
	r.Get("/:job", uiJobSingle)
}
