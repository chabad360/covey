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

	"github.com/chabad360/covey/host/types"
	"golang.org/x/crypto/ssh"
)

// SSHHost implements host.Host and host.HostInfo
type SSHHost struct {
	types.HostInfo
	PrivateKey []byte
	PublicKey  []byte
	HostKey    []byte
	Username   string
	Port       string
	config     *ssh.ClientConfig
}

type newHostInfo struct {
	Server   string
	Port     string
	Username string
	Password string
	Name     string
	Plugin   string
}

type plugin struct{}

// Plugin is the plugin for managing SSH based hosts.
var Plugin plugin

// Run runs a command on the server.
func (h *SSHHost) Run(args []string) (*bytes.Buffer, error) {
	var b bytes.Buffer
	// Create an initial connection
	client, err := ssh.Dial("tcp", h.Server+":"+h.Port, h.config)
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

// NewHost creates an SSHHost
func (p *plugin) NewHost(hostJSON []byte) (types.Host, error) {
	h := SSHHost{}

	var hostInfo newHostInfo
	if err := json.Unmarshal(hostJSON, &hostInfo); err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: hostInfo.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(hostInfo.Password),
		},
		HostKeyCallback: hostKeyCallback(&h),
	}

	// Create an initial connection
	client, err := ssh.Dial("tcp", hostInfo.Server+":"+hostInfo.Port, config)
	if err != nil {
		return nil, err
	}

	output, err := sshRun(client, "/usr/bin/whoami")
	if err != nil {
		return nil, err
	}
	// Verify that everything has gone right
	output = output[0 : len(output)-1]
	if string(output) != hostInfo.Username {
		return nil, fmt.Errorf("%v is not %v", string(output), hostInfo.Username)
	}
	log.Println("Successfully logged into server")
	// Generate SSH Keys add add the public key to the authorized_keys file.
	err = generateKeysAndAddKeys(&h)
	if err != nil {
		return nil, err
	}
	// log.Println("Generated SSH keys")
	if _, err := sshRun(client, fmt.Sprint("echo '", string(h.PublicKey), "' | tee -a .ssh/authorized_keys")); err != nil {
		return nil, err
	}
	client.Close()

	h.Name = hostInfo.Name
	h.Username = hostInfo.Username
	h.Server = hostInfo.Server
	h.Port = hostInfo.Port
	h.Plugin = "ssh"

	if err := hostFactory(&h); err != nil {
		return nil, err
	}

	return &h, nil
}

func (p *plugin) LoadHost(hostJSON []byte) (types.Host, error) {
	var h SSHHost
	if err := json.Unmarshal(hostJSON, &h); err != nil {
		return nil, err
	}
	log.Println("Loading", h.Name)

	if err := hostFactory(&h); err != nil {
		return nil, err
	}

	return &h, nil
}

func generateKeysAndAddKeys(h *SSHHost) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(privateKey),
	}
	h.PrivateKey = pem.EncodeToMemory(&privBlock)

	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}
	h.PublicKey = ssh.MarshalAuthorizedKey(publicKey)

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

func hostKeyCallback(h *SSHHost) ssh.HostKeyCallback {
	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		if len(h.HostKey) <= 0 {
			h.HostKey = key.Marshal()
		}
		return nil
	}
}

func hostFactory(h *SSHHost) error {
	signer, err := ssh.ParsePrivateKey(h.PrivateKey)
	if err != nil {
		return err
	}

	hostKey, err := ssh.ParsePublicKey(h.HostKey)
	if err != nil {
		return err
	}
	// log.Printf("loaded host key")

	config := &ssh.ClientConfig{
		User: h.Username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	// Create an initial connection
	client, err := ssh.Dial("tcp", h.Server+":"+h.Port, config)
	if err != nil {
		return err
	}
	output, err := sshRun(client, "/usr/bin/whoami")
	if err != nil {
		return err
	}
	// Verify that everything has gone right
	output = output[0 : len(output)-1]
	if string(output) != h.Username {
		return fmt.Errorf("%v is not %v", string(output), h.Username)
	}
	client.Close()
	h.config = config

	// log.Printf("Created Host")
	return nil
}
