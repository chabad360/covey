package task

import (
	"fmt"
	"github.com/chabad360/covey/models"
	"net/http"
	"strings"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/storage"
	"github.com/chabad360/covey/ui"
	"github.com/go-playground/pure/v5"
)

func uiTasks(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	var tasks []models.Task
	result := storage.DB.Find(&tasks)
	ui.ErrorWriter(w, result.Error)

	p := &ui.Page{
		Title:   "Tasks",
		URL:     strings.Split(r.URL.Path, "/"),
		Details: struct{ Tasks []models.Task }{Tasks: tasks},
	}

	t := ui.GetTemplate("tasksAll")
	err := t.ExecuteTemplate(w, "base", p)
	ui.ErrorWriter(w, err)
}

func uiTaskSingle(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	vars := pure.RequestVars(r)

	task, ok := storage.GetTask(vars.URLParam("task"))
	ui.ErrorWriter404(w, vars.URLParam("task"), ok)

	p := &ui.Page{
		Title:   fmt.Sprintf("Task %s", vars.URLParam("task")),
		URL:     strings.Split(r.URL.Path, "/"),
		Details: struct{ Task *models.Task }{Task: task},
	}

	t := ui.GetTemplate("tasksSingle")
	err := t.ExecuteTemplate(w, "base", p)
	ui.ErrorWriter(w, err)
}

// UITaskNew returns the form for creating a new task.
func UITaskNew(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	var nodes []string
	storage.DB.Table("nodes").Select("name").Scan(&nodes)

	p := &ui.Page{
		Title: "New Task",
		URL:   strings.Split(r.URL.Path, "/"),
		Details: struct {
			Plugins []string
			Nodes   []string
		}{[]string{"shell"}, nodes},
	}

	t := ui.GetTemplate("tasksNew")
	err := t.ExecuteTemplate(w, "base", p)
	ui.ErrorWriter(w, err)
}

// RegisterUIHandlers registers the HTTP handlers for the task UI.
func RegisterUIHandlers(r pure.IRouteGroup) {
	r.Get("", uiTasks)
	r.Get("/:task", uiTaskSingle)
}
