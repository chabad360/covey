package plugin

import (
	"net/http"

	"github.com/go-playground/pure/v5"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/ui"
)

func GetPlugin(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	vars := pure.RequestVars(r)

	plugin, err := GetTaskPlugin(vars.URLParam("plugin"))
	ui.ErrorWriter(w, err)

	t := ui.GetTemplate("form")
	err = t.ExecuteTemplate(w, "form", plugin.GetInputs())
	ui.ErrorWriter(w, err)
}
