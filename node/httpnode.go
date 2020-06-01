package node

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/node/types"
	"github.com/gorilla/mux"
)

var (
	nodes      = make(map[string]types.INode)
	nodesShort = make(map[string]string)
	nodesName  = make(map[string]string)
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

	nodes[t.GetID()] = t
	nodesShort[t.GetIDShort()] = t.GetID()
	nodesName[t.GetName()] = t.GetID()
	j, err := json.MarshalIndent(nodes, "", "  ")
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}
	f, err := os.Create("./config/nodes.json")
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}
	defer f.Close()
	if err = f.Chmod(0600); err != nil {
		common.ErrorWriter(w, err)
		return
	}
	if _, err = f.Write(j); err != nil {
		common.ErrorWriter(w, err)
		return
	}

	j, err = json.MarshalIndent(t, "", "  ")
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
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 %s not found", vars["node"])
		return
	}
	w.Header().Add("Content-Type", "application/json")

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
	j := new(struct {
		Result []string
	})
	j.Result = []string{}
	var c []byte
	for _, byte := range b.Bytes() {
		if byte != '\n' {
			c = append(c, byte)
		} else {
			j.Result = append(j.Result, string(c))
			c = nil
		}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(j)
}

// NodeGet returns a JSON representation of the specified node, GET /api/v1/node/{node}
func nodeGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	n, ok := GetNode(vars["node"])
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 %s not found", vars["node"])
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
