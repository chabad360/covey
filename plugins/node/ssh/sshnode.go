package main

import (
	"bytes"
	"log"
	"strings"

	"golang.org/x/crypto/ssh"
)

// GetDetails returns the details of the node.
func (n *Node) GetDetails() interface{} { return n.Details }

// Run runs a command on the server.
func (n *Node) Run(args []string) (*bytes.Buffer, chan int, error) {
	var b bytes.Buffer
	// Create an initial connection
	client, err := ssh.Dial("tcp", n.Details.Server+":"+n.Details.Port, n.Details.config)
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
