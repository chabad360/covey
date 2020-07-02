package node

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/chabad360/covey/common"
	nodeSSH "github.com/chabad360/covey/node/ssh"
	"github.com/chabad360/covey/node/types"
	"github.com/go-playground/pure/v5"
	json "github.com/json-iterator/go"
)

// nodeNew adds a new node using the specified plugin.
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

	n, err := nodeSSH.NewNode(reqBody)
	common.ErrorWriter(w, err)

	if err = addNode(n); err != nil {
		common.ErrorWriter(w, err)
	}

	w.Header().Add("Location", "/api/v1/node/"+n.GetID())
	common.Write(w, n)
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
func RegisterHandlers(newRoute pure.IRouteGroup, singleRoute pure.IRouteGroup) {
	log.Println("Registering Node module API handlers...")
	newRoute.Post("/add", nodeNew)

	n := singleRoute.Group("/:node")
	n.Get("", nodeGet)
}
