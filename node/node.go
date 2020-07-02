package node

import (
	"log"

	"github.com/chabad360/covey/node/types"
)

// GetNode checks if a node with the identifier exists and returns it.
func GetNode(identifier string) (*types.Node, bool) {
	n, pk, puk, hk, err := getNodeAndKeys(identifier)
	if err != nil {
		log.Printf("GetNode: %v", err)
		return nil, false
	}
	t, err := loadNode(n, pk, puk, hk)
	if err != nil {
		log.Printf("LoadNode: %v", err)
		return nil, false
	}

	return t, true
}
