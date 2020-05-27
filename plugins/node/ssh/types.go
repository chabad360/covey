package main

import (
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
