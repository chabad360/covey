package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/node/types"
	"github.com/chabad360/covey/storage"
	"github.com/gorilla/mux"
)

// NodeNew adds a new node using the specified plugin.
func nodeNew(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var node types.Node
	reqBody, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(reqBody, &node); err != nil {
		common.ErrorWriter(w, err)
		return
	}

	if _, ok := GetNode(node.Name); ok {
		common.ErrorWriter(w, fmt.Errorf("Duplicate node: %v", node.Name))
		return
	}

	p, err := loadPlugin(node.Plugin)
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}

	t, err := p.NewNode(reqBody)
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}

	if err = storage.AddItem("nodes", t.GetID(), t.GetIDShort(), t); err != nil {
		common.ErrorWriter(w, err)
		return
	}
	log.Println("Stored Node")

	z, err := storage.GetItem("nodes", t.GetID(), t)
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}
	j, err := json.MarshalIndent(z, "", "  ")
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}

	w.Header().Set("Location", "/api/v1/node/"+t.GetName())
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(j))
}

// NodeRun runs a command the specified node, POST /api/v1/node/{node}
func nodeRun(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	n, ok := GetNode(vars["node"])
	w.Header().Add("Content-Type", "application/json")
	if !ok {
		common.ErrorWriter(w, errors.New("404 not found"))
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
	b, _, err := n.Run(s.Cmd)
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}
	z := new(struct {
		Result []string
	})
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
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(z)
}

// NodeGet returns a JSON representation of the specified node, GET /api/v1/node/{node}
func nodeGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	n, ok := GetNode(vars["node"])
	if !ok {
		common.ErrorWriter(w, errors.New("404 not found"))
		return
	}
	w.Header().Add("Content-Type", "application/json")

	j, err := json.MarshalIndent(n, "", "\t")
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(j))

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
