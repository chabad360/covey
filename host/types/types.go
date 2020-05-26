package types

import "bytes"

// HostPlugin defines what a host plugin should look like
type HostPlugin interface {
	// NewHost returns a new host
	NewHost(hostJSON []byte) (Host, error)

	// LoadHost loads the json representation of each host host
	LoadHost(hostJSON []byte) (Host, error)
}

// HostInfo contains information about a host and must be implemented alongside the Host interface.
type HostInfo struct {
	Name   string
	Server string
	Plugin string
}

// Host defines the generic host
type Host interface {
	// Run a command on the host
	Run(args []string) (*bytes.Buffer, error)

	// GetName returns the name of the host
	GetName() string
}
