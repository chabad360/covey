package nodeSSH

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"

	"fmt"
	"net"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/node/types"
	json "github.com/json-iterator/go"
	"golang.org/x/crypto/ssh"
)

func generateAndAddKeys(n *types.Node) error {
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

func hostKeyCallback(n *types.Node) ssh.HostKeyCallback {
	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		if len(n.HostKey) <= 0 {
			n.HostKey = key.Marshal()
		}
		return nil
	}
}

func nodeFactory(n *types.Node) error {
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
	client, err := ssh.Dial("tcp", n.IP+":"+n.Port, config)
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
	n.Config = config

	// log.Printf("Created Node")
	return nil
}

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
	err = generateAndAddKeys(node)
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
