package safe

// TaskPluginInterface defines the interface used by task plugins.
type TaskPluginInterface interface {
	// GetCommand returns the command to run the server.
	GetCommand(Task) (string, error)

	// GetFetchCommand returns a command to run which will be used to fetch relevant information about the node, and a callback that returns JSON metadata to send the output too.
	// GetFetchCommand() (string, func([]string) ([]byte, error)) TODO: add support for metadata

	// GetInputs returns the inputs that the plugin takes.
	// GetInputs([]byte) ([]byte, error) TODO: add support for customizing inputs based on metadata
	GetInputs() Form
}
