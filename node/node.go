package node

import (
	"log"

	nodeSSH "github.com/chabad360/covey/node/ssh"
	"github.com/chabad360/covey/node/types"
	"github.com/chabad360/covey/storage"
)

// GetNode checks if a node with the identifier exists and returns it.
func GetNode(identifier string) (*types.Node, bool) {
	n, err := storage.GetItem("nodes", identifier)
	if err != nil {
		log.Printf("GetNode: %v", err)
		return nil, false
	}
	t, err := nodeSSH.LoadNode(n)
	if err != nil {
		log.Printf("LoadNode: %v", err)
		return nil, false
	}

	return t, true
}
