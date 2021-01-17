//go:generate yaegi extract -license ../LICENSE -name plugins github.com/chabad360/covey/models/safe

package plugins

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/chabad360/plugins"

	"github.com/chabad360/covey/config"
	"github.com/chabad360/covey/models/safe"
	"github.com/chabad360/covey/plugins/task/shell"
)

var (
	Host    *plugins.PluginHost
	Symbols = make(map[string]map[string]reflect.Value)
)

func Init() error {
	Host, _ = plugins.NewPluginHost(config.Config.Plugins.PluginsFolder, config.Config.Plugins.PluginsCacheFolder, Symbols)
	Host.AddPluginType("task", (*safe.TaskPluginInterface)(nil))

	Host.AddInternalPlugin(reflect.ValueOf(shell.GetPlugin()), plugins.PluginConfig{
		Name:        "Shell",
		Description: "Runs a shell command on a node.",
		PluginType:  "task",
	})

	if err := Host.LoadPlugins(); err != nil && !errors.Is(err, plugins.ErrNoDirectorySpecified) {
		return err
	}

	return nil
}

func GetTaskPlugin(pluginName string) (safe.TaskPluginInterface, error) {
	p, ok := Host.GetPlugin(pluginName)
	if !ok {
		return nil, fmt.Errorf("Failed to load plugin: %v", pluginName)
	}

	return p.Interface().(safe.TaskPluginInterface), nil
}
