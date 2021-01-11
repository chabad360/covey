package node

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-playground/pure/v5"
	json "github.com/json-iterator/go"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/models"
	"github.com/chabad360/covey/storage"
)

// nodeNew adds a new node.
func nodeNew(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	var n models.Node

	reqBody, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(reqBody, &n); err != nil {
		common.ErrorWriterCustom(w, err, http.StatusBadRequest)
	}

	if _, ok := storage.GetNode(n.Name); ok {
		common.ErrorWriterCustom(w, fmt.Errorf("duplicate node: %v", n.Name), http.StatusConflict)
	}

	err := newNode(&n)
	common.ErrorWriter(w, err)

	err = storage.AddNode(&n)
	common.ErrorWriter(w, err)

	w.Header().Add("Location", "/api/v1/node/"+n.ID)
	common.Write(w, n)
}

func nodesGet(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	var q storage.Query
	err := q.Setup(r)
	common.ErrorWriter(w, err)

	var nodes interface{}

	if q.Expand {
		var n []models.Node
		err = q.Query("nodes", &n)
		nodes = n
	} else {
		var n []string
		err = q.Query("nodes", &n)
		nodes = n
	}
	common.ErrorWriter(w, err)

	common.Write(w, nodes)
}

// nodeGet returns a JSON representation of the specified node, GET /api/v1/node/{node}.
func nodeGet(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	vars := pure.RequestVars(r)
	n, ok := storage.GetNode(vars.URLParam("node"))
	common.ErrorWriter404(w, vars.URLParam("node"), ok)

	common.Write(w, n)
}

func nodeDelete(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	vars := pure.RequestVars(r)
	n, ok := storage.GetNode(vars.URLParam("node"))
	common.ErrorWriter404(w, vars.URLParam("node"), ok)

	err := storage.DeleteNode(n)
	common.ErrorWriter(w, err)

	common.Write(w, vars.URLParam("node"))
}

// RegisterHandlers adds the handlers for the node module.
func RegisterHandlers(r pure.IRouteGroup) {
	log.Println("Registering Node module API handlers...")
	r.Post("", nodeNew)
	r.Get("", nodesGet)

	n := r.Group("/:node")
	n.Get("", nodeGet)
	n.Delete("", nodeDelete)
}
