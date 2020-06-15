package ui

import (
	"log"
	"net/http"

	"html/template"

	"github.com/chabad360/covey/common"
	"github.com/go-playground/pure/v5"
)

var (
// baseTemplate = template.New("base")
)

func dashboard(w http.ResponseWriter, _ *http.Request) {
	str, ok := common.FS.String("/base.html")
	if !ok {
		log.Fatal("Missing files")
	}
	t, err := template.New("base").Parse(str)
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}
	str, ok = common.FS.String("/sidebar.html")
	if !ok {
		log.Fatal("Missing files")
	}
	t.New("sidebar").Parse(str)
	err = t.ExecuteTemplate(w, "base", struct {
		Title string
		Body  string
	}{Title: "Dashboard", Body: "{{ .Title }}"})
	if err != nil {
		common.ErrorWriter(w, err)
	}
}

// RegisterHandlers registers the handlers for the ui module.
func RegisterHandlers(r pure.IRouteGroup) {
	r.Get("/dashboard", dashboard)
}
