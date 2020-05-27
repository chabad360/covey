package types

import (
	"bytes"
	"encoding/hex"
)

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
	ID      string      `json:"id,omitempty"`
}

// GetID returns the ID of the task.
func (n *Node) GetID() string { return n.ID }

// GetIDShort returns the first 8 bytes of the task ID.
func (n *Node) GetIDShort() string { x, _ := hex.DecodeString(n.ID); return hex.EncodeToString(x[:8]) }

// INode defines the generic node interface.
type INode interface {
	// Run a command on the node
	Run(args []string) (*bytes.Buffer, chan int, error)

	// GetName returns the name of the node
	GetName() string
}
