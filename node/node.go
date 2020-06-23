package node

import (
	"log"

	nodeSSH "github.com/chabad360/covey/node/ssh"
	"github.com/chabad360/covey/node/types"
	"github.com/chabad360/covey/storage"
)

// LoadConfig loads up the stored nodes
// func LoadConfig() {
// 	log.Println("Placeholder")
// }

func loadNode(nodeJSON []byte) (*types.Node, error) {
	t, err := nodeSSH.LoadNode(nodeJSON)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// GetNode checks if a node with the identifier exists and returns it.
func GetNode(identifier string) (*types.Node, bool) {
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
