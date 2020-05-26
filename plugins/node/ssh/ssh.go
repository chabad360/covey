package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/chabad360/covey/node/types"
	"golang.org/x/crypto/ssh"
)

// SSHNode implements node.Node and node.NodeInfo
type SSHNode struct {
	types.NodeInfo
	PrivateKey []byte
	PublicKey  []byte
	HostKey    []byte
	Username   string
	Port       string
	config     *ssh.ClientConfig
}

type newNodeInfo struct {
	Server   string
	Port     string
	Username string
	Password string
	Name     string
	Plugin   string
}

type plugin struct{}

// Plugin is the plugin for managing SSH based nodes.
var Plugin plugin

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

func generateKeysAndAddKeys(n *SSHNode) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(privateKey),
	}
	n.PrivateKey = pem.EncodeToMemory(&privBlock)

	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}
	n.PublicKey = ssh.MarshalAuthorizedKey(publicKey)

	return nil
}

func sshRun(client *ssh.Client, command string) ([]byte, error) {
	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	output, err := session.Output(command)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func hostKeyCallback(n *SSHNode) ssh.HostKeyCallback {
	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		if len(n.HostKey) <= 0 {
			n.HostKey = key.Marshal()
		}
		return nil
	}
}

func nodeFactory(n *SSHNode) error {
	signer, err := ssh.ParsePrivateKey(n.PrivateKey)
	if err != nil {
		return err
	}

	hostKey, err := ssh.ParsePublicKey(n.HostKey)
	if err != nil {
		return err
	}
	// log.Printf("loaded node key")

	config := &ssh.ClientConfig{
		User: n.Username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	// Create an initial connection
	client, err := ssh.Dial("tcp", n.Server+":"+n.Port, config)
	if err != nil {
		return err
	}
	output, err := sshRun(client, "/usr/bin/whoami")
	if err != nil {
		return err
	}
	// Verify that everything has gone right
	output = output[0 : len(output)-1]
	if string(output) != n.Username {
		return fmt.Errorf("%v is not %v", string(output), n.Username)
	}
	client.Close()
	n.config = config

	// log.Printf("Created Node")
	return nil
}
