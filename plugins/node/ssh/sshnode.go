package main

import (
	"bytes"
	"strings"

	"golang.org/x/crypto/ssh"
)

// GetName returns the name of the Node
func (n *Node) GetName() string {
	return n.Name
}

// Run runs a command on the server.
func (n *Node) Run(args []string) (*bytes.Buffer, error) {
	var b bytes.Buffer
	// Create an initial connection
	client, err := ssh.Dial("tcp", n.Details.Server+":"+n.Details.Port, n.Details.config)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	session.Stdout = &b
	go func() error {
		if err := session.Run(strings.Join(args, " ")); err != nil {
			return err
		}
		return nil
	}()
	return &b, nil
}
