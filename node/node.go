package node

import (
	"fmt"
	"io"
	"log"
	"os"
	"plugin"

	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/chabad360/covey/node/types"
	"github.com/gorilla/mux"
)

var (
	nodes = make(map[string]types.Node)
)

func loadPlugin(pluginName string) (types.NodePlugin, error) {
	p, err := plugin.Open("./plugins/node/" + pluginName + ".so")
	if err != nil {
		return nil, err
	}

	n, err := p.Lookup("Plugin")
	if err != nil {
		return nil, err
	}

	var s types.NodePlugin
	s, ok := n.(types.NodePlugin)
	if !ok {
		return nil, fmt.Errorf(pluginName, " does not provide a NodePlugin")
	}

	return s, nil
}

// NewNode adds a new node using the specified plugin.
func NewNode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var node types.NodeInfo
	reqBody, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(reqBody, &node); err != nil {
		errorWriter(w, err)
		return
	}

	for h := range nodes {
		if h == node.Name {
			errorWriter(w, fmt.Errorf("Duplicate node: %v", node.Name))
			return
		}
	}

	p, err := loadPlugin(node.Plugin)
	if err != nil {
		errorWriter(w, err)
		return
	}

	h, err := p.NewNode(reqBody)
	if err != nil {
		errorWriter(w, err)
		return
	}

	nodes[h.GetName()] = h
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

	j, err = json.MarshalIndent(h, "", "  ")
	if err != nil {
		errorWriter(w, err)
		return
	}
	fmt.Fprintf(w, string(j))
}

// RegisterHandlers adds the mux handlers for the node module.
func RegisterHandlers(r *mux.Router) {
	r.HandleFunc("/add", NewNode).Methods("POST")
}

// LoadConfig loads up the stored nodes
func LoadConfig() {
	log.Println("Loading Node Config")
	f, err := os.Open("./config/nodes.json")
	if err != nil {
		log.Println("Error loading node config")
		return
	}
	defer f.Close()

	var h map[string]json.RawMessage
	if err = json.NewDecoder(f).Decode(&h); err != nil {
		log.Fatal(err)
	}

	// Make this dynamic
	var plugins = make(map[string]types.NodePlugin)
	p, err := loadPlugin("ssh")
	if err != nil {
		log.Fatal(err)
	}
	plugins["ssh"] = p

	for _, node := range h {
		var z types.NodeInfo
		j, err := node.MarshalJSON()
		if err != nil {
			log.Fatal(err)
		}
		if err := json.Unmarshal(j, &z); err != nil {
			log.Fatal(err)
		}

		t, err := plugins[z.Plugin].LoadNode(j)
		if err != nil {
			log.Fatal(err)
		}
		nodes[t.GetName()] = t

		r, err := t.Run([]string{"echo", "Hello World"})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(r)

	}
}

func errorWriter(w io.Writer, err error) {
	fmt.Fprintf(w, "{'error':'%s'}", err)
}
