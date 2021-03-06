package job

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/pure/v5"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/models"
	"github.com/chabad360/covey/storage"
	"github.com/chabad360/covey/ui"
)

func uiJobs(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	var jobs []models.Job
	result := storage.DB.Find(&jobs)
	ui.ErrorWriter(w, result.Error)

	p := &ui.Page{
		Title:   "Jobs",
		URL:     strings.Split(r.URL.Path, "/"),
		Details: struct{ Jobs []models.Job }{Jobs: jobs},
	}

	t := ui.GetTemplate("jobsAll")
	err := t.ExecuteTemplate(w, "base", p)
	ui.ErrorWriter(w, err)
}

func uiJobSingle(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	vars := pure.RequestVars(r)

	job, ok := storage.GetJobWithFullHistory(vars.URLParam("job"))
	ui.ErrorWriter404(w, vars.URLParam("job"), ok)

	if r.URL.Query().Get("run") == "true" { // TODO: ?
		j, _ := storage.GetJob(vars.URLParam("job"))
		_, err := run(j)
		ui.ErrorWriter(w, err)

		job, _ = storage.GetJobWithFullHistory(vars.URLParam("job"))
	}

	p := &ui.Page{
		Title:   fmt.Sprintf("Job %s", vars.URLParam("job")),
		URL:     strings.Split(r.URL.Path, "/"),
		Details: struct{ Job *models.JobWithTasks }{job},
	}

	t := ui.GetTemplate("jobsSingle")
	err := t.ExecuteTemplate(w, "base", p)
	ui.ErrorWriter(w, err)
}

// UIJobNew returns the form for creating a new task.
func UIJobNew(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	var nodes []string
	result := storage.DB.Table("nodes").Select("name").Scan(&nodes)
	ui.ErrorWriter(w, result.Error)

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
	ui.ErrorWriter(w, err)
}

// RegisterUIHandlers registers the HTTP handlers for the jobs UI.
func RegisterUIHandlers(r pure.IRouteGroup) {
	r.Get("", uiJobs)
	r.Get("/:job", uiJobSingle)
}
