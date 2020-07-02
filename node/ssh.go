package node

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/chabad360/covey/asset"
	"log"

	"fmt"
	"net"

	scp "github.com/bramvdbogaerde/go-scp"
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

func hostKeyCallback(n *types.Node) ssh.HostKeyCallback {
	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		if len(n.HostKey) == 0 {
			n.HostKey = key.Marshal()
		}
		return nil
	}
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

func installAgent(node string, config *ssh.ClientConfig) error {
	client := scp.NewClient(node, config)
	if err := client.Connect(); err != nil {
		return err
	}

	f, err := asset.FS.Open("/agent")
	if err != nil {
		return err
	}

	defer client.Close()
	defer f.Close()

	if err = client.CopyFile(f, "/tmp/agent", "0755"); err != nil {
		return err
	}

	return nil
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
	output = output[:len(output)-1]
	if string(output) != n.Username {
		return fmt.Errorf("%v is not %v", string(output), n.Username)
	}

	client.Close()

	n.Config = config
	// log.Printf("Created Node")
	return nil
}

func newNode(nodeJSON []byte) (*types.Node, error) {
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
	output = output[:len(output)-1]
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
	if _, err = sshRun(client, fmt.Sprint("echo '", string(node.PublicKey),
		"' | tee -a .ssh/authorized_keys")); err != nil {
		return nil, err
	}

	if err = nodeFactory(node); err != nil {
		return nil, err
	}
	node.ID = common.GenerateID(node)

	if err = installAgent(node.IP+":"+node.Port, config); err != nil {
		return nil, err
	}

	if output, err := sshRun(client, "sudo mv /tmp/agent /usr/bin/"); err != nil {
		return nil, fmt.Errorf("mv /tmp/agent: %v, error: %v", output, err)
	}

	if _, err := sshRun(client, fmt.Sprintf(`sudo mkdir /etc/covey && echo 'AGENT_ID="%s"
AGENT_HOST="%s"' | sudo tee /etc/covey/agent.conf`, node.ID, "192.168.56.1")); err != nil {
		return nil, err
	} // Add config file for agent
	client.Close()

	return node, nil
}

func loadNode(nodeJSON []byte, privateKey []byte, publicKey []byte, hostKey []byte) (*types.Node, error) {
	var n types.Node
	if err := json.Unmarshal(nodeJSON, &n); err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}
	n.PrivateKey = privateKey
	n.PublicKey = publicKey
	n.HostKey = hostKey

	if err := nodeFactory(&n); err != nil {
		return nil, err
	}

	return &n, nil
}
