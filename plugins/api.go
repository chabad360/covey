package plugins

import (
	"log"
	"net/http"

	"github.com/go-playground/pure/v5"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/models/safe"
	"github.com/chabad360/covey/storage"
)

func pluginsGet(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	var q storage.Query
	err := q.Setup(r)
	common.ErrorWriter(w, err)

	var plugins []interface{}
	var i, l int

	for k, v := range Host.Plugins {
		if i < q.Offset {
			i++
			continue
		}
		if l >= q.Limit {
			break
		}
		if q.Expand {
			plugins = append(plugins, v)
		} else {
			plugins = append(plugins, k)
		}
		i++
		l++
	}
	common.ErrorWriter(w, err)

	common.Write(w, plugins)
}

func pluginGet(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	vars := pure.RequestVars(r)

	p, ok := Host.Plugins[vars.URLParam("plugin")]
	common.ErrorWriter404(w, vars.URLParam("plugin"), ok)

	common.Write(w, p)
}

func inputsGet(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	vars := pure.RequestVars(r)

	pI, ok := Host.GetPlugin(vars.URLParam("plugin"))
	p, ok := pI.Interface().(safe.TaskPluginInterface)
	common.ErrorWriter404(w, vars.URLParam("plugin"), ok)

	common.Write(w, p.GetInputs())
}

// RegisterHandlers registers the handlers for the plugin module.
func RegisterHandlers(r pure.IRouteGroup) {
	log.Println("Registering Plugin module API handlers...")

	r.Get("", pluginsGet)

	t := r.Group("/:plugin")
	t.Get("", pluginGet)
	t.Get("/inputs", inputsGet)
}
