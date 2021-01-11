package node

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net"

	"github.com/bramvdbogaerde/go-scp"
	"golang.org/x/crypto/ssh"

	"github.com/chabad360/covey/asset"
	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/config"
	"github.com/chabad360/covey/models"
)

func generateAndAddKeys(n *models.Node) error {
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

func hostKeyCallback(n *models.Node) ssh.HostKeyCallback {
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

func installAgent(node string, id string, cfg *ssh.ClientConfig, sshClient *ssh.Client) error {
	log.Println("installing agent...")

	client := scp.NewClient(node, cfg)

	f, err := asset.FS.Open("/agent/agent")
	if err != nil {
		return fmt.Errorf("open /agent/agent: %v", err)
	}
	defer f.Close()

	f2, err := asset.FS.Open("/agent/covey-agent.service")
	if err != nil {
		return fmt.Errorf("open /agent/covey-agent.service: %v", err)
	}
	defer f2.Close()

	if err := client.Connect(); err != nil {
		return err
	}
	defer client.Close()
	if err = client.CopyFile(f, "/tmp/agent", "0755"); err != nil {
		return fmt.Errorf("copy /agent/agent: %v", err)
	}

	if err := client.Connect(); err != nil {
		return err
	}
	defer client.Close()
	if err = client.CopyFile(f2, "/tmp/covey-agent.service", "0644"); err != nil {
		return fmt.Errorf("copy /agent/covey-agent.service: %v", err)
	}
	log.Println("Copied files")

	if _, err = sshRun(sshClient, `sudo chown root:root /tmp/agent \
		&& sudo chown root:root /tmp/covey-agent.service`); err != nil {
		return fmt.Errorf("chown: %v", err)
	}
	if _, err = sshRun(sshClient, "sudo mv /tmp/agent /usr/bin/"); err != nil {
		return fmt.Errorf("install agent: %v", err)
	}
	if _, err = sshRun(sshClient, "sudo mv /tmp/covey-agent.service /usr/lib/systemd/system/"); err != nil {
		return fmt.Errorf("install service: %v", err)
	}

	if _, err := sshRun(sshClient, fmt.Sprintf(`sudo mkdir /etc/covey; echo 'AGENT_ID="%s"
AGENT_HOST="%s"
AGENT_PORT="%s"' | sudo tee /etc/covey/agent.conf`,
		id, config.Config.Daemon.Host, config.Config.Daemon.Port)); err != nil {
		return fmt.Errorf("install config: %v", err)
	} // Add config file for agent
	if _, err = sshRun(sshClient, "sudo systemctl enable --now covey-agent.service"); err != nil {
		return fmt.Errorf("install service: %v", err)
	}

	return nil
}

func nodeFactory(n *models.Node) error {
	err := n.Setup()
	if err != nil {
		return err
	}

	// Create an initial connection
	client, err := ssh.Dial("tcp", n.IP+":"+n.Port, n.Config)
	if err != nil {
		return err
	}

	defer client.Close()

	output, err := sshRun(client, "/usr/bin/whoami")
	if err != nil {
		return err
	}

	// Verify that everything has gone right
	output = output[:len(output)-1]
	if string(output) != n.Username {
		return fmt.Errorf("%v is not %v", string(output), n.Username)
	}

	return nil
}

func newNode(node *models.Node) error {
	cfg := &ssh.ClientConfig{
		User: node.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(node.Password),
		},
		HostKeyCallback: hostKeyCallback(node),
	}

	// Create an initial connection
	client, err := ssh.Dial("tcp", node.IP+":"+node.Port, cfg)
	if err != nil {
		return err
	}

	defer client.Close()

	output, err := sshRun(client, "/usr/bin/whoami")
	if err != nil {
		return err
	}

	// Verify that we can run commands and get what we expected.
	output = output[:len(output)-1]
	if string(output) != node.Username {
		return fmt.Errorf("newNode: %v is not %v", string(output), node.Username)
	}
	log.Println("Successfully logged into node")

	// Generate SSH Keys add add the public key to the authorized_keys file.
	err = generateAndAddKeys(node)
	if err != nil {
		return err
	}
	// log.Println("Generated SSH keys")
	if _, err = sshRun(client, fmt.Sprint("echo '", string(node.PublicKey),
		"' | tee -a .ssh/authorized_keys")); err != nil {
		return err
	}

	if err = nodeFactory(node); err != nil {
		return err
	}
	node.ID = common.GenerateID(node)

	return installAgent(node.IP+":"+node.Port, node.ID, cfg, client)
}
