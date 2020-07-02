package job

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/job/types"
	"github.com/chabad360/covey/storage"
	"github.com/chabad360/covey/ui"
	"github.com/go-playground/pure/v5"
)

func uiJobs(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	var jobs []types.Job
	err := storage.DB.QueryRow(context.Background(),
		"SELECT jsonb_agg(to_jsonb(jobs) - 'log' - 'details') FROM jobs").Scan(&jobs)
	common.ErrorWriter(w, err)

	p := &ui.Page{
		Title:   "Jobs",
		URL:     strings.Split(r.URL.Path, "/"),
		Details: struct{ Jobs []types.Job }{Jobs: jobs},
	}

	t := ui.GetTemplate("jobsAll")
	err = t.ExecuteTemplate(w, "base", p)
	common.ErrorWriter(w, err)
}

func uiJobSingle(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	vars := pure.RequestVars(r)

	job, ok := GetJobWithTasks(vars.URLParam("job")) // This goes first so we dont need to confirm existence anymore.
	if !ok {
		common.ErrorWriter404(w, vars.URLParam("job"))
	}

	if r.URL.Query().Get("run") == "true" {
		j, _ := GetJob(vars.URLParam("job"))
		j.Run()

		err := UpdateJob(*j)
		common.ErrorWriter(w, err)

		job, _ = GetJobWithTasks(vars.URLParam("job"))
	}

	p := &ui.Page{
		Title:   fmt.Sprintf("Job %s", vars.URLParam("job")),
		URL:     strings.Split(r.URL.Path, "/"),
		Details: struct{ Job *types.JobWithTasks }{job},
	}

	t := ui.GetTemplate("jobsSingle")
	err := t.ExecuteTemplate(w, "base", p)
	common.ErrorWriter(w, err)
}

// UIJobNew returns the form for creating a new task.
func UIJobNew(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	var nodes []string
	err := storage.DB.QueryRow(context.Background(), "SELECT jsonb_agg(name) FROM nodes;").Scan(&nodes)
	common.ErrorWriter(w, err)

	p := &ui.Page{
		Title: "New Job",
		URL:   strings.Split(r.URL.Path, "/"),
		Details: struct {
			Plugins []string
			Nodes   []string
		}{[]string{"shell"}, nodes},
	}

	t := ui.GetTemplate("jobsNew")
	err = t.ExecuteTemplate(w, "base", p)
	common.ErrorWriter(w, err)
}

// RegisterUIHandlers registers the HTTP handlers for the jobs UI.
func RegisterUIHandlers(r pure.IRouteGroup) {
	r.Get("", uiJobs)
	r.Get("/:job", uiJobSingle)
}
