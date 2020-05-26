package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/chabad360/covey/host"
	"golang.org/x/crypto/ssh"
)

// SSHHost implements host.Host and host.HostInfo
type SSHHost struct {
	host.HostInfo
	client     *ssh.Client
	privateKey []byte
	publicKey  []byte
	hostKey    []byte
	username   string
}

type plugin struct{}

// Plugin is the plugin for managing SSH based hosts.
var Plugin plugin

// Run runs a command on the server.
func (h *SSHHost) Run(args []string) (*bytes.Buffer, error) {
	var b bytes.Buffer

	session, err := h.client.NewSession()
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

// Plugin returns the plugin associated with the Host.
func (h *SSHHost) Plugin() string {
	return "ssh"
}

// NewHost creates an SSHHost
func (p *plugin) NewHost(NewHostInfo *host.NewHostInfo) (host.Host, error) {
	h := SSHHost{}

	config := &ssh.ClientConfig{
		User: NewHostInfo.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(NewHostInfo.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Create an initial connection
	client, err := ssh.Dial("tcp", NewHostInfo.Server+":"+NewHostInfo.Port, config)
	if err != nil {
		return nil, err
	}
	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	output, err := session.Output("/usr/bin/whoami")
	if err != nil {
		return nil, err
	}
	// Verify that everything has gone right
	if string(output) != NewHostInfo.Username {
		return nil, fmt.Errorf("%v is not %v", output, NewHostInfo.Username)
	}

	// Generate SSH Keys add add the public key to the authorized_keys file.
	err = generateKeysAndAddKeys(&h)
	if err != nil {
		return nil, err
	}
	if err := session.Run(fmt.Sprint("echo '", h.publicKey, "' | tee -a .ssh/authorized_keys")); err != nil {
		return nil, err
	}
	output, err = session.Output("cat /etc/ssh/ssh_host_rsa_key.pub")
	if err != nil {
		return nil, err
	}
	h.hostKey = output
	session.Close()

	signer, err := ssh.ParsePrivateKey(h.privateKey)
	if err != nil {
		return nil, err
	}

	hostKey, err := ssh.ParsePublicKey(h.hostKey)
	if err != nil {
		return nil, err
	}

	config = &ssh.ClientConfig{
		User: NewHostInfo.Username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	// Create an initial connection
	client, err = ssh.Dial("tcp", NewHostInfo.Server+":"+NewHostInfo.Port, config)
	if err != nil {
		return nil, err
	}
	session, err = client.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	output, err = session.Output("/usr/bin/whoami")
	if err != nil {
		return nil, err
	}
	// Verify that everything has gone right
	if string(output) != h.username {
		return nil, fmt.Errorf("%v is not %v", output, h.username)
	}
	h.client = client

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
	h.privateKey = pem.EncodeToMemory(&privBlock)

	publicKey, err := ssh.NewPublicKey(privateKey)
	if err != nil {
		return nil
	}
	h.publicKey = ssh.MarshalAuthorizedKey(publicKey)

	return nil
}
