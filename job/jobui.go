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
	var jobs []byte
	err := storage.DB.QueryRow(context.Background(),
		"SELECT jsonb_agg(to_jsonb(jobs) - 'log' - 'details') FROM jobs").Scan(&jobs)
	if err != nil {
		common.ErrorWriter(w, err)
	}
	p := &ui.Page{
		Title:   "Jobs",
		URL:     strings.Split(r.URL.Path, "/"),
		Details: struct{ Jobs string }{Jobs: string(jobs)},
	}
	t := ui.GetTemplate("jobsAll")
	err = t.ExecuteTemplate(w, "base", p)
	if err != nil {
		common.ErrorWriter(w, err)
	}
}

func uiJobSingle(w http.ResponseWriter, r *http.Request) {
	vars := pure.RequestVars(r)
	job, ok := GetJobWithTasks(vars.URLParam("job"))
	if !ok {
		common.ErrorWriter404(w, vars.URLParam("job"))
	}

	p := &ui.Page{
		Title:   fmt.Sprintf("Job %s", vars.URLParam("job")),
		URL:     strings.Split(r.URL.Path, "/"),
		Details: struct{ Job *types.JobWithTasks }{job},
	}

	t := ui.GetTemplate("jobsSingle")
	err := t.ExecuteTemplate(w, "base", p)
	if err != nil {
		common.ErrorWriter(w, err)
	}
}

// RegisterUIHandlers registers the HTTP handlers for the jobs UI.
func RegisterUIHandlers(r pure.IRouteGroup) {
	r.Get("", uiJobs)
	r.Get("/:job", uiJobSingle)
}
