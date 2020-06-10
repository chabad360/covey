package main

import (
	"fmt"
	"log"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/node/types"
	json "github.com/json-iterator/go"
	"golang.org/x/crypto/ssh"
)

// Plugin is the plugin for managing SSH based nodes.
var Plugin plugin

// NewNode creates an SSHNode
func (p *plugin) NewNode(nodeJSON []byte) (types.INode, error) {
	var nodeInfo newNodeInfo
	if err := json.Unmarshal(nodeJSON, &nodeInfo); err != nil {
		return nil, err
	}

	x := Node{
		Details: &SSHNode{
			Username: nodeInfo.Username,
			Server:   nodeInfo.Server,
			Port:     nodeInfo.Port,
		},
	}
	x.Plugin = "ssh"
	x.Name = nodeInfo.Name

	config := &ssh.ClientConfig{
		User: nodeInfo.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(nodeInfo.Password),
		},
		HostKeyCallback: hostKeyCallback(x.Details),
	}

	// Create an initial connection
	client, err := ssh.Dial("tcp", x.Details.Server+":"+x.Details.Port, config)
	if err != nil {
		return nil, err
	}

	output, err := sshRun(client, "/usr/bin/whoami")
	if err != nil {
		return nil, err
	}
	// Verify that we can run commands and get what we expected.
	output = output[0 : len(output)-1]
	if string(output) != nodeInfo.Username {
		return nil, fmt.Errorf("%v is not %v", string(output), nodeInfo.Username)
	}
	log.Println("Successfully logged into node")
	// Generate SSH Keys add add the public key to the authorized_keys file.
	err = generateKeysAndAddKeys(x.Details)
	if err != nil {
		return nil, err
	}
	// log.Println("Generated SSH keys")
	if _, err := sshRun(client,
		fmt.Sprint("echo '", string(x.Details.PublicKey), "' | tee -a .ssh/authorized_keys")); err != nil {
		return nil, err
	}
	client.Close()

	if err := nodeFactory(x.Details); err != nil {
		return nil, err
	}
	x.ID = common.GenerateID(x)

	return &x, nil
}

// LoadNode takes the node information and convertes it into a usable node.
func (p *plugin) LoadNode(nodeJSON []byte) (types.INode, error) {
	var n Node
	if err := json.Unmarshal(nodeJSON, &n); err != nil {
		return nil, err
	}

	if err := nodeFactory(n.Details); err != nil {
		return nil, err
	}

	return &n, nil
}
