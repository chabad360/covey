package node

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/node/types"
	"github.com/chabad360/covey/storage"
	taskTypes "github.com/chabad360/covey/task/types"
	"github.com/chabad360/covey/ui"
	"github.com/go-playground/pure/v5"
)

func uiNodes(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	var nodes []types.Node
	err := storage.DB.QueryRow(context.Background(),
		"SELECT jsonb_agg(to_jsonb(nodes) - 'private_key' - 'public_key' - 'host_key') FROM nodes;").Scan(&nodes)
	common.ErrorWriter(w, err)

	p := &ui.Page{
		Title:   "Nodes",
		URL:     strings.Split(r.URL.Path, "/"),
		Details: struct{ Nodes []types.Node }{nodes},
	}
	t := ui.GetTemplate("nodesAll")
	err = t.ExecuteTemplate(w, "base", p)
	common.ErrorWriter(w, err)
}

func uiNodeSingle(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	vars := pure.RequestVars(r)

	node, ok := GetNode(vars.URLParam("node"))
	if !ok {
		common.ErrorWriter404(w, vars.URLParam("node"))
	}
	var tasks []taskTypes.Task

	err := storage.DB.QueryRow(context.Background(),
		"SELECT jsonb_agg(to_jsonb(tasks)) FROM tasks WHERE node = $1;", node.Name).Scan(&tasks)
	common.ErrorWriter(w, err)

	p := &ui.Page{
		Title: fmt.Sprintf("Node %s", vars.URLParam("node")),
		URL:   strings.Split(r.URL.Path, "/"),
		Details: struct {
			Node  *types.Node
			Tasks []taskTypes.Task
			Host  string
		}{node, tasks, "localhost"},
	}

	t := ui.GetTemplate("nodesSingle")
	err = t.ExecuteTemplate(w, "base", p)
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
