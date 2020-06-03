package node

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/node/types"
	"github.com/gorilla/mux"
)

// nodeNew adds a new node using the specified plugin.
func nodeNew(w http.ResponseWriter, r *http.Request) {
	var node types.Node
	reqBody, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(reqBody, &node); err != nil {
		common.ErrorWriterCustom(w, err, http.StatusBadRequest)
		return
	}

	if _, ok := GetNode(node.Name); ok {
		common.ErrorWriterCustom(w, fmt.Errorf("Duplicate node: %v", node.Name), http.StatusConflict)
		return
	}

	plugin, err := loadPlugin(node.Plugin)
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}

	n, err := plugin.NewNode(reqBody)
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}

	if err = addNode(n); err != nil {
		common.ErrorWriter(w, err)
		return
	}

	w.Header().Add("Location", "/api/v1/node/"+n.GetID())
	common.Write(w, n)

}

// nodeRun runs a command the specified node, POST /api/v1/node/{node}
func nodeRun(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	n, ok := GetNode(vars["node"])
	if !ok {
		common.ErrorWriter404(w, vars["node"])
		return
	}

	var s struct {
		Cmd []string
	}
	reqBody, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(reqBody, &s); err != nil {
		common.ErrorWriter(w, err)
		return
	}
	if len(s.Cmd) == 0 {
		common.ErrorWriter(w, fmt.Errorf("Missing command"))
		return
	}

	b, d, err := n.Run(s.Cmd)
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}

	z := new(struct {
		Result   []string `json:"result"`
		ExitCode int      `json:"exit_code"`
	})
	e := <-d
	c := []byte{}
	l := []string{}
	for _, bb := range b.Bytes() {
		if bb == '\n' {
			l = append(l, string(c))
			c = nil
		} else {
			c = append(c, bb)
		}
	}
	z.Result = l
	z.ExitCode = e
	common.Write(w, z)
}

// nodeGet returns a JSON representation of the specified node, GET /api/v1/node/{node}
func nodeGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	n, ok := GetNode(vars["node"])
	if !ok {
		common.ErrorWriter404(w, vars["node"])
		return
	}

	common.Write(w, n)

}

// RegisterHandlers adds the mux handlers for the node module.
func RegisterHandlers(r *mux.Router) {
	log.Println("Registering Node module API handlers...")

	r.HandleFunc("/new", nodeNew).Methods("POST")
	r.HandleFunc("/{node}", nodeRun).Methods("POST")
	r.HandleFunc("/{node}", nodeGet).Methods("GET")

	err := r.Walk(common.Walk)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()
}
