package node

import (
	"fmt"
	"github.com/chabad360/covey/models"
	"github.com/chabad360/covey/storage"
	"net/http"
	"strings"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/ui"
	"github.com/go-playground/pure/v5"
)

func uiNodes(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	var nodes []models.Node
	result := storage.DB.Find(&nodes)
	common.ErrorWriter(w, result.Error)

	p := &ui.Page{
		Title:   "Nodes",
		URL:     strings.Split(r.URL.Path, "/"),
		Details: struct{ Nodes []models.Node }{nodes},
	}
	t := ui.GetTemplate("nodesAll")
	err := t.ExecuteTemplate(w, "base", p)
	common.ErrorWriter(w, err)
}

func uiNodeSingle(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	vars := pure.RequestVars(r)
	node, ok := storage.GetNode(vars.URLParam("node"))
	common.ErrorWriter404(w, vars.URLParam("node"), ok)

	var tasks []models.Task
	result := storage.DB.Table("tasks").Where("node = ?", node.ID).Or("node = ?", node.Name).Find(&tasks)
	common.ErrorWriter(w, result.Error)

	p := &ui.Page{
		Title: fmt.Sprintf("Node %s", vars.URLParam("node")),
		URL:   strings.Split(r.URL.Path, "/"),
		Details: struct {
			Node  *models.Node
			Tasks []models.Task
			Host  string
		}{node, tasks, "localhost"},
	}

	t := ui.GetTemplate("nodesSingle")
	err := t.ExecuteTemplate(w, "base", p)
	common.ErrorWriter(w, err)
}

// UINodeNew returns the form for creating a new task.
func UINodeNew(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	p := &ui.Page{
		Title: "New Node",
		URL:   strings.Split(r.URL.Path, "/"),
	}

	t := ui.GetTemplate("nodesNew")
	err := t.ExecuteTemplate(w, "base", p)
	common.ErrorWriter(w, err)
}

// RegisterUIHandlers registers the HTTP handlers for the nodes UI.
func RegisterUIHandlers(r pure.IRouteGroup) {
	r.Get("", uiNodes)
	r.Get("/:node", uiNodeSingle)
}
