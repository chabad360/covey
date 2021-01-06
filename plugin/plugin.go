package plugin

import (
	"fmt"
	"github.com/chabad360/covey/config"
	"github.com/chabad360/covey/models/safe"
	"github.com/chabad360/plugins"
)

var (
	Host *plugins.PluginHost
)

func Init() error {
	Host = plugins.NewPluginHost(config.Config.Plugins.PluginsFolder, config.Config.Plugins.PluginsCacheFolder)
	Host.AddPluginType("task", (*safe.TaskPluginInterface)(nil))
	return Host.LoadPlugins()
}

func GetTaskPlugin(pluginName string) (safe.TaskPluginInterface, error) {
	p, err := Host.GetPlugin(pluginName)
	if err != nil {
		return nil, fmt.Errorf("Failed to load plugin: %v", pluginName)
	}

	return p.(safe.TaskPluginInterface), nil
}
