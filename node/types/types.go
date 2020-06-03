package types

import (
	"bytes"
	"encoding/hex"
)

// NodePlugin defines what a node plugin should look like
type NodePlugin interface {
	// NewNode returns a new node
	NewNode(nodeJSON []byte) (INode, error)

	// LoadNode loads the json representation of each node.
	LoadNode(nodeJSON []byte) (INode, error)
}

// INode defines the generic node interface.
type INode interface {
	// Run a command on the node
	Run(args []string) (*bytes.Buffer, chan int, error)

	// GetName returns the name of the node
	GetName() string

	// GetID returns the id of the node.
	GetID() string

	// GetIDShort returns the first 8 bytes of the node ID.
	GetIDShort() string

	// GetPlugin returns the plugin of the node.
	GetPlugin() string

	// GetDetails returns the details of the node.
	GetDetails() interface{}
}

// Node contains information about a node and must be implemented alongside the INode interface.
type Node struct {
	Name    string      `json:"name"`
	Plugin  string      `json:"plugin"`
	Details interface{} `json:"details"`
	ID      string      `json:"id"`
}

// GetName returns the name of the Node.
func (n *Node) GetName() string { return n.Name }

// GetID returns the ID of the node.
func (n *Node) GetID() string { return n.ID }

// GetIDShort returns the first 8 bytes of the node ID.
func (n *Node) GetIDShort() string { x, _ := hex.DecodeString(n.ID); return hex.EncodeToString(x[:8]) }

// GetPlugin returns the plugin of the node.
func (n *Node) GetPlugin() string { return n.Plugin }
