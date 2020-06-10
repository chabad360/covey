package node

import (
	"fmt"
	"log"
	"plugin"

	json "github.com/json-iterator/go"

	"github.com/chabad360/covey/node/types"
	"github.com/chabad360/covey/storage"
)

// LoadConfig loads up the stored nodes
// func LoadConfig() {
// 	log.Println("Placeholder")
// }

func loadPlugin(pluginName string) (types.NodePlugin, error) {
	p, err := plugin.Open("./plugins/node/" + pluginName + ".so")
	if err != nil {
		return nil, err
	}

	n, err := p.Lookup("Plugin")
	if err != nil {
		return nil, err
	}

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
	n, err := storage.GetItem("nodes", identifier)
	if err != nil {
		log.Println(err)
		return nil, false
	}
	x, err := loadNode(n)
	if err != nil {
		log.Println(err)
		return nil, false
	}

	return x, true
}
