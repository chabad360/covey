package nodeSSH

import (
	"fmt"
	"log"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/node/types"
	json "github.com/json-iterator/go"
	"golang.org/x/crypto/ssh"
)

// NewNode creates an types.Node
func NewNode(nodeJSON []byte) (*types.Node, error) {
	var node *types.Node
	if err := json.Unmarshal(nodeJSON, &node); err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: node.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(node.Password),
		},
		HostKeyCallback: hostKeyCallback(node),
	}

	// Create an initial connection
	client, err := ssh.Dial("tcp", node.IP+":"+node.Port, config)
	if err != nil {
		return nil, err
	}

	output, err := sshRun(client, "/usr/bin/whoami")
	if err != nil {
		return nil, err
	}
	// Verify that we can run commands and get what we expected.
	output = output[0 : len(output)-1]
	if string(output) != node.Username {
		return nil, fmt.Errorf("%v is not %v", string(output), node.Username)
	}
	log.Println("Successfully logged into node")
	// Generate SSH Keys add add the public key to the authorized_keys file.
	err = generateKeysAndAddKeys(node)
	if err != nil {
		return nil, err
	}
	// log.Println("Generated SSH keys")
	if _, err := sshRun(client,
		fmt.Sprint("echo '", string(node.PublicKey), "' | tee -a .ssh/authorized_keys")); err != nil {
		return nil, err
	}
	client.Close()

	if err := nodeFactory(node); err != nil {
		return nil, err
	}
	node.ID = common.GenerateID(node)

	return node, nil
}

// LoadNode takes the node information and converts it into a usable node.
func LoadNode(nodeJSON []byte) (*types.Node, error) {
	var n types.Node
	if err := json.Unmarshal(nodeJSON, &n); err != nil {
		return nil, err
	}

	if err := nodeFactory(&n); err != nil {
		return nil, err
	}

	return &n, nil
}
