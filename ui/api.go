package ui

import (
	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/plugin"
	"github.com/go-playground/pure/v5"
	json "github.com/json-iterator/go"
	"net/http"
)

func pluginInputs(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()
	vars := pure.RequestVars(r)
	p, err := plugin.GetTaskPlugin(vars.URLParam("plugin"))
	ErrorWriter(w, err)

	formBytes, err := p.GetInputs()
	ErrorWriter(w, err)

	var form Form
	ErrorWriter(w, json.Unmarshal(formBytes, &form))

	t := GetTemplate("pluginForm")
	ErrorWriter(w, t.ExecuteTemplate(w, "form", form))
}
