package ui

import (
	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/plugin"
	"github.com/go-playground/pure/v5"
	"net/http"
)

func pluginInputs(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()
	vars := pure.RequestVars(r)
	p, err := plugin.GetTaskPlugin(vars.URLParam("plugin"))
	ErrorWriter(w, err)

	form := p.GetInputs()

	t := GetTemplate("pluginForm")
	ErrorWriter(w, t.ExecuteTemplate(w, "form", form))
}
