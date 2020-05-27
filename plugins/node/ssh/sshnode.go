package main

import (
	"bytes"
	"strings"

	"golang.org/x/crypto/ssh"
)

// GetName returns the name of the Node
func (n *SSHNode) GetName() string {
	return n.Name
}

// Run runs a command on the server.
func (n *SSHNode) Run(args []string) (*bytes.Buffer, error) {
	var b bytes.Buffer
	// Create an initial connection
	client, err := ssh.Dial("tcp", n.Server+":"+n.Port, n.config)
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
	if err := session.Run(strings.Join(args, " ")); err != nil {
		return nil, err
	}
	return &b, nil
}
