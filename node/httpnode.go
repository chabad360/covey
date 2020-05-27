package node

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/chabad360/covey/node/types"
	"github.com/gorilla/mux"
)

var (
	nodes = make(map[string]types.Node)
)

// NewNode adds a new node using the specified plugin.
func NewNode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var node types.NodeInfo
	reqBody, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(reqBody, &node); err != nil {
		errorWriter(w, err)
		return
	}

	if _, ok := nodes[node.Name]; ok {
		errorWriter(w, fmt.Errorf("Duplicate node: %v", node.Name))
		return
	}

	p, err := loadPlugin(node.Plugin)
	if err != nil {
		errorWriter(w, err)
		return
	}

	t, err := p.NewNode(reqBody)
	if err != nil {
		errorWriter(w, err)
		return
	}

	nodes[t.GetName()] = t
	j, err := json.MarshalIndent(nodes, "", "  ")
	if err != nil {
		errorWriter(w, err)
		return
	}
	f, err := os.Create("./config/nodes.json")
	if err != nil {
		errorWriter(w, err)
		return
	}
	defer f.Close()
	if err = f.Chmod(0600); err != nil {
		errorWriter(w, err)
		return
	}
	if _, err = f.Write(j); err != nil {
		errorWriter(w, err)
		return
	}

	j, err = json.MarshalIndent(t, "", "  ")
	if err != nil {
		errorWriter(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Location", "/api/v1/nodes/"+t.GetName())
	fmt.Fprintf(w, string(j))
}

// RunNode runs a command the specified node, POST /api/v1/node/{node}
func RunNode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	n, ok := nodes[vars["node"]]
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
		errorWriter(w, err)
		return
	}
	if len(s.Cmd) == 0 {
		errorWriter(w, fmt.Errorf("Missing command"))
	}

	b, err := n.Run(s.Cmd)
	if err != nil {
		errorWriter(w, err)
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

// GetNode returns a JSON representation of the specified node, GET /api/v1/node/{node}
func GetNode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	n, ok := nodes[vars["node"]]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 %s not found", vars["node"])
		return
	}
	w.Header().Add("Content-Type", "application/json")

	j, err := json.MarshalIndent(n, "", "  ")
	if err != nil {
		errorWriter(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(j))

}

// RegisterHandlers adds the mux handlers for the node module.
func RegisterHandlers(r *mux.Router) {
	log.Println("Registering Node module API handlers...")
	r.HandleFunc("/new", NewNode).Methods("POST")
	r.HandleFunc("/{node}", RunNode).Methods("POST")
	r.HandleFunc("/{node}", GetNode).Methods("GET")
}

func errorWriter(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "{'error':'%s'}", err)
}
