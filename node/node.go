package node

import (
	"fmt"
	"log"
	"os"
	"plugin"

	"encoding/json"

	"github.com/chabad360/covey/node/types"
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
	}
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
