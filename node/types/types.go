package types

import (
	"bytes"
	"encoding/hex"
	"log"
	"strings"

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

// Run runs a command on the server.
func (n *Node) Run(args []string) (*bytes.Buffer, chan int, error) {
	var b bytes.Buffer
	// Create an initial connection
	client, err := ssh.Dial("tcp", n.IP+":"+n.Port, n.Config)
	if err != nil {
		return nil, nil, err
	}
	session, err := client.NewSession()
	if err != nil {
		return nil, nil, err
	}
	c := make(chan int)

	session.Stdout = &b
	// session.Stderr = &b
	go func() {
		if err := session.Run(strings.Join(args, " ")); err != nil {
			if msg, ok := err.(*ssh.ExitError); ok {
				log.Println(msg.ExitStatus())
				c <- msg.ExitStatus()
			}
		} else {
			c <- 0
		}
		close(c)
		session.Close()
		client.Close()
	}()
	return &b, c, nil
}
