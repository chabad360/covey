package models

import (
	"encoding/hex"
	"gorm.io/gorm"
	"time"

	"golang.org/x/crypto/ssh"
)

// Node contains information about a node.
type Node struct {
	Name       string            `json:"name" gorm:"<-:create;notnull;unique"`
	ID         string            `json:"id" gorm:"<-:create;primarykey"`
	IDShort    string            `json:"-" gorm:"<-:create;notnull;unique"`
	PrivateKey []byte            `json:"-" gorm:"<-:create;notnull"`
	PublicKey  []byte            `json:"-" gorm:"<-:create;notnull"`
	HostKey    []byte            `json:"-" gorm:"<-:create;notnull;unique"`
	IP         string            `json:"ip" gorm:"<-:create;notnull"`
	Username   string            `json:"username" gorm:"<-:create;notnull;default:root"`
	Password   string            `json:"password,omitempty" gorm:"-"`
	Port       string            `json:"port" gorm:"<-:create;default:22"`
	Config     *ssh.ClientConfig `json:"-" gorm:"-"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// GetIDShort returns the first 8 bytes of the node ID.
func (n *Node) GetIDShort() string { x, _ := hex.DecodeString(n.ID); return hex.EncodeToString(x[:8]) }

// BeforeCreate gets the short ID before saving.
func (n *Node) BeforeCreate(_ *gorm.DB) (err error) {
	n.IDShort = n.GetIDShort()
	return nil
}

// AfterFind runs the setup for a node
func (n *Node) AfterFind(_ *gorm.DB) (err error) {
	return n.Setup()
}

// Setup is responsible for creating an SSH client.
func (n *Node) Setup() error {
	signer, err := ssh.ParsePrivateKey(n.PrivateKey)
	if err != nil {
		return err
	}

	hostKey, err := ssh.ParsePublicKey(n.HostKey)
	if err != nil {
		return err
	}

	n.Config = &ssh.ClientConfig{
		User: n.Username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	return nil
}
