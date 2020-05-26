package types

import "bytes"

// NodePlugin defines what a node plugin should look like
type NodePlugin interface {
	// NewNode returns a new node
	NewNode(nodeJSON []byte) (Node, error)

	// LoadNode loads the json representation of each node node
	LoadNode(nodeJSON []byte) (Node, error)
}

// NodeInfo contains information about a node and must be implemented alongside the Node interface.
type NodeInfo struct {
	Name   string
	Server string
	Plugin string
}

// Node defines the generic node
type Node interface {
	// Run a command on the node
	Run(args []string) (*bytes.Buffer, error)

	// GetName returns the name of the node
	GetName() string
}
