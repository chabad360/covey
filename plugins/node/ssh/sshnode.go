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
func (n *Node) Run(args []string) (*bytes.Buffer, chan int, error) {
	var b bytes.Buffer
	// Create an initial connection
	client, err := ssh.Dial("tcp", n.Details.Server+":"+n.Details.Port, n.Details.config)
	if err != nil {
		return nil, nil, err
	}
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		return nil, nil, err
	}
	c := make(chan int)

	session.Stdout = &b
	go func() {
		session.Start(strings.Join(args, " "))
		if err := session.Wait(); err != nil {
			if msg, ok := err.(*ssh.ExitError); ok {
				c <- msg.ExitStatus()
			}
		} else {
			c <- 0
		}
		close(c)
	}()
	return &b, c, nil
}
