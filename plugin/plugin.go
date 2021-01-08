package plugin

import (
	"fmt"
	"github.com/chabad360/covey/config"
	"github.com/chabad360/covey/models/safe"
	"github.com/chabad360/plugins"
	"reflect"
)

var (
	Host    *plugins.PluginHost
	Symbols = make(map[string]map[string]reflect.Value)
)

func Init() error {
	Host = plugins.NewPluginHost(config.Config.Plugins.PluginsFolder, config.Config.Plugins.PluginsCacheFolder, Symbols)
	Host.AddPluginType("task", (*safe.TaskPluginInterface)(nil))
	return Host.LoadPlugins()
}

func GetTaskPlugin(pluginName string) (safe.TaskPluginInterface, error) {
	p, ok := Host.GetPlugin(pluginName)
	if !ok {
		return nil, fmt.Errorf("Failed to load plugin: %v", pluginName)
	}

	return p.Interface().(safe.TaskPluginInterface), nil
}
