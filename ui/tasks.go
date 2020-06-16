package ui

import (
	"context"
	"net/http"
	"strings"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/storage"
)

func tasks(w http.ResponseWriter, r *http.Request) {
	var tasks []byte
	err := storage.DB.QueryRow(context.Background(),
		"SELECT jsonb_agg(to_jsonb(tasks) - 'id' - 'log' - 'details') FROM tasks").Scan(&tasks)
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
