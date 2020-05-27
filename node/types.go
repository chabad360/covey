package node

import "bytes"

// NodePlugin defines what a node plugin should look like
type NodePlugin interface {
	// NewNode returns a new node
	NewNode(nodeJSON []byte) (INode, error)

	// LoadNode loads the json representation of each node node
	LoadNode(nodeJSON []byte) (INode, error)
}

// Node contains information about a node and must be implemented alongside the INode interface.
type Node struct {
	Name    string      `json:"name,omitempty"`
	Plugin  string      `json:"plugin,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

// INode defines the generic node interface.
type INode interface {
	// Run a command on the node
	Run(args []string) (*bytes.Buffer, error)

	// GetName returns the name of the node
	GetName() string
}
