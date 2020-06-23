package types

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/ssh"
)

// Node contains information about a node.
type Node struct {
	Name       string            `json:"name"`
	ID         string            `json:"id"`
	PrivateKey []byte            `json:"private_key"`
	PublicKey  []byte            `json:"public_key"`
	HostKey    []byte            `json:"host_key"`
	IP         string            `json:"ip"`
	Username   string            `json:"username"`
	Password   string            `json:"password,omitempty"`
	Port       string            `json:"port"`
	Config     *ssh.ClientConfig `json:"-"`
}

// GetName returns the name of the Node.
func (n *Node) GetName() string { return n.Name }

// GetID returns the ID of the node.
func (n *Node) GetID() string { return n.ID }

// GetIDShort returns the first 8 bytes of the node ID.
func (n *Node) GetIDShort() string { x, _ := hex.DecodeString(n.ID); return hex.EncodeToString(x[:8]) }

// Run is a stub implementation of the Run method.
func (n *Node) Run(_ []string) (*bytes.Buffer, chan int, error) {
	return nil, nil, fmt.Errorf("not implemented")
}
