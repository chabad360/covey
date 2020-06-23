package nodeSSH

import (
	"bytes"
	"log"
	"strings"

	"github.com/chabad360/covey/node/types"
	"golang.org/x/crypto/ssh"
)

// Run runs a command on the server.
func (n *types.Node) Run(args []string) (*bytes.Buffer, chan int, error) {
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
