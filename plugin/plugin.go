package plugin

import (
	"github.com/chabad360/covey/config"
	"github.com/chabad360/covey/models"
	"github.com/chabad360/plugins"
)

type TaskPluginInterface interface {
	// GetCommand returns the command to run the server.
	GetCommand(task models.Task) (string, error)

	// GetFetchCommand returns a command to run which will be used to fetch relevant information about the node, and a callback that returns JSON metadata to send the output too.
	// GetFetchCommand() (string, func([]string) ([]byte, error)) TODO: add support for metadata

	// GetInputs takes the JSON metadata and returns a JSON that can be converted to a ui.Form object.
	// GetInputs([]byte) ([]byte, error) TODO: add support for customizing inputs based on metadata
	GetInputs() ([]byte, error)
}

var (
	Host *plugins.PluginHost
)

func Init() error {
	Host = plugins.NewPluginHost(config.Config.Plugins.PluginsFolder, config.Config.Plugins.PluginsCacheFolder)
	Host.AddPluginType("task", (*TaskPluginInterface)(nil))
	return Host.LoadPlugins()
}

func GetTaskPlugin(pluginName string) (TaskPluginInterface, error) {
	p, err := Host.GetPlugin(pluginName)
	if err != nil {
		return nil, err
	}

	return p.(TaskPluginInterface), nil
}
