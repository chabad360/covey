package main

import (
	"github.com/chabad360/covey/node/types"
	"golang.org/x/crypto/ssh"
)

// SSHNode contains the details of an SSH node
type SSHNode struct {
	PrivateKey []byte `json:"private_key,omitempty"`
	PublicKey  []byte `json:"public_key,omitempty"`
	HostKey    []byte `json:"host_key,omitempty"`
	Server     string `json:"server,omitempty"`
	Username   string `json:"username,omitempty"`
	Port       string `json:"port,omitempty"`
	config     *ssh.ClientConfig
}

// Node is a generic Node type
type Node struct {
	types.Node
	Details *SSHNode `json:"details,omitempty"`
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
