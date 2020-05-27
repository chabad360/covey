package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/chabad360/covey/node/types"
	"golang.org/x/crypto/ssh"
)

// Plugin is the plugin for managing SSH based nodes.
var Plugin plugin

// NewNode creates an SSHNode
func (p *plugin) NewNode(nodeJSON []byte) (types.Node, error) {
	n := SSHNode{}

	var nodeInfo newNodeInfo
	if err := json.Unmarshal(nodeJSON, &nodeInfo); err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: nodeInfo.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(nodeInfo.Password),
		},
		HostKeyCallback: hostKeyCallback(&n),
	}

	// Create an initial connection
	client, err := ssh.Dial("tcp", nodeInfo.Server+":"+nodeInfo.Port, config)
	if err != nil {
		return nil, err
	}

	output, err := sshRun(client, "/usr/bin/whoami")
	if err != nil {
		return nil, err
	}
	// Verify that everything has gone right
	output = output[0 : len(output)-1]
	if string(output) != nodeInfo.Username {
		return nil, fmt.Errorf("%v is not %v", string(output), nodeInfo.Username)
	}
	log.Println("Successfully logged into server")
	// Generate SSH Keys add add the public key to the authorized_keys file.
	err = generateKeysAndAddKeys(&n)
	if err != nil {
		return nil, err
	}
	// log.Println("Generated SSH keys")
	if _, err := sshRun(client, fmt.Sprint("echo '", string(n.PublicKey), "' | tee -a .ssh/authorized_keys")); err != nil {
		return nil, err
	}
	client.Close()

	n.Name = nodeInfo.Name
	n.Username = nodeInfo.Username
	n.Server = nodeInfo.Server
	n.Port = nodeInfo.Port
	n.Plugin = "ssh"

	if err := nodeFactory(&n); err != nil {
		return nil, err
	}

	return &n, nil
}

func (p *plugin) LoadNode(nodeJSON []byte) (types.Node, error) {
	var n SSHNode
	if err := json.Unmarshal(nodeJSON, &n); err != nil {
		return nil, err
	}
	log.Println("Loading", n.Name)

	if err := nodeFactory(&n); err != nil {
		return nil, err
	}

	return &n, nil
}
