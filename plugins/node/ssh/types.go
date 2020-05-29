package main

import (
	"github.com/chabad360/covey/node/types"
	"golang.org/x/crypto/ssh"
)

// SSHNode contains the details of an SSH node
type SSHNode struct {
	PrivateKey []byte `json:"private_key"`
	PublicKey  []byte `json:"public_key"`
	HostKey    []byte `json:"host_key"`
	Server     string `json:"server"`
	Username   string `json:"username"`
	Port       string `json:"port"`
	config     *ssh.ClientConfig
}

// Node is a generic Node type
type Node struct {
	types.Node
	Details *SSHNode `json:"details"`
}

type newNodeInfo struct {
	Server   string `json:"server"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Plugin   string `json:"plugin"`
}

type plugin struct{}
