package ui

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/storage"
	"github.com/chabad360/covey/task/types"
	"github.com/go-playground/pure/v5"
	json "github.com/json-iterator/go"
)

func tasks(w http.ResponseWriter, r *http.Request) {
	var tasks []byte
	err := storage.DB.QueryRow(context.Background(),
		"SELECT jsonb_agg(to_jsonb(tasks) - 'log' - 'details') FROM tasks").Scan(&tasks)
	if err != nil {
		common.ErrorWriter(w, err)
	}
	p := &page{
		Title:   "Tasks",
		URL:     strings.Split(r.URL.Path, "/"),
		Details: struct{ Tasks string }{Tasks: string(tasks)},
	}
	t := getTemplate("tasksAll")
	err = t.ExecuteTemplate(w, "base", p)
	if err != nil {
		common.ErrorWriter(w, err)
	}
}

func task(w http.ResponseWriter, r *http.Request) {
	vars := pure.RequestVars(r)
	var j []byte
	err := storage.DB.QueryRow(context.Background(),
		"SELECT to_jsonb(tasks) FROM tasks WHERE id = $1", vars.URLParam("id")).Scan(&j)
	if err != nil {
		common.ErrorWriter(w, err)
	}
	var task types.Task
	err = json.Unmarshal(j, &task)
	if err != nil {
		common.ErrorWriter(w, err)
	}

	p := &page{
		Title:   fmt.Sprintf("Task %s", vars.URLParam("id")),
		URL:     strings.Split(r.URL.Path, "/"),
		Details: struct{ Task types.Task }{Task: task},
	}

	t := getTemplate("tasksSingle")
	err = t.ExecuteTemplate(w, "base", p)
	if err != nil {
		common.ErrorWriter(w, err)
	}
}
