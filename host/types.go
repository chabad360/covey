package host

import "bytes"

// HostPlugin defines what a host plugin should look like
type HostPlugin interface {
	// NewHost returns a new host
	NewHost(newHostInfo *NewHostInfo) (Host, error)
}

// NewHostInfo contains the info about a new host and is passed to the specified plugin
type NewHostInfo struct {
	Server   string
	Port     string
	Username string
	Password string
	Name     string
	Plugin   string
}

// HostInfo contains information about a host and must be implemented alongside the Host interface.
type HostInfo struct {
	Name   string
	Server string
}

// Host defines the generic host
type Host interface {
	// Run a command on the host
	Run(args []string) (*bytes.Buffer, error)

	// Returns the name of the plugin that handles it
	Plugin() string
}
