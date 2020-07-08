package node

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/node/types"
	"github.com/go-playground/pure/v5"
	json "github.com/json-iterator/go"
)

// nodeNew adds a new node.
func nodeNew(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	var node types.Node

	reqBody, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(reqBody, &node); err != nil {
		common.ErrorWriterCustom(w, err, http.StatusBadRequest)
	}

	if _, ok := GetNode(node.Name); ok {
		common.ErrorWriterCustom(w,
			fmt.Errorf("duplicate node: %v", node.Name), http.StatusConflict)
	}

	n, err := newNode(reqBody)
	common.ErrorWriter(w, err)

	if err = addNode(n); err != nil {
		common.ErrorWriter(w, err)
	}

	w.Header().Add("Location", "/api/v1/node/"+n.GetID())
	common.Write(w, n)
}

func nodesGet(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()
	refreshDB()

	var nodes []string
	err := db.QueryRow(context.Background(), "SELECT jsonb_agg(name) FROM nodes;").Scan(&nodes)
	common.ErrorWriter(w, err)

	common.Write(w, nodes)
}

// nodeGet returns a JSON representation of the specified node, GET /api/v1/node/{node}.
func nodeGet(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	vars := pure.RequestVars(r)
	n, ok := GetNode(vars.URLParam("node"))
	if !ok {
		common.ErrorWriter404(w, vars.URLParam("node"))
	}

	common.Write(w, n)
}

// RegisterHandlers adds the handlers for the node module.
func RegisterHandlers(r pure.IRouteGroup) {
	log.Println("Registering Node module API handlers...")
	r.Post("", nodeNew)
	r.Get("", nodesGet)
	n := r.Group("/:node")
	n.Get("", nodeGet)
}
