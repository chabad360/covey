package node

import (
	"fmt"
	"log"
	"os"
	"plugin"

	"encoding/json"

	"github.com/chabad360/covey/node/types"
	"github.com/chabad360/covey/storage"
)

// LoadConfig loads up the stored nodes
func LoadConfig() {
	log.Println("Loading Node Config")
	f, err := os.Open("./config/nodes.json")
	if err != nil {
		log.Println("Error loading node config")
		return
	}
	defer f.Close()
}

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

func loadNode(nodeJSON []byte) (types.INode, error) {
	var z types.Node
	if err := json.Unmarshal(nodeJSON, &z); err != nil {
		log.Fatal(err)
	}
	p, err := loadPlugin(z.Plugin)
	if err != nil {
		return nil, err
	}
	t, err := p.LoadNode(nodeJSON)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// GetNode checks if a node with the identifier exists and returns it.
func GetNode(identifier string) (types.INode, bool) {
	var t *types.Node
	n, err := storage.GetItem("nodes", identifier, t)
	if err != nil {
		log.Println(err)
		return nil, false
	}
	j, err := json.Marshal(n)
	if err != nil {
		log.Println(err)
		return nil, false
	}
	x, err := loadNode(j)
	if err != nil {
		log.Println(err)
		return nil, false
	}

	return x, true
}
