package ui

import (
	"net/http"

	"html/template"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/ui/templates"
	"github.com/gorilla/mux"
)

var (
// baseTemplate = template.New("base")
)

func dashboard(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("base").Parse(templates.Base)
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}
	t.New("sidebar").Parse(templates.Sidebar)
	err = t.ExecuteTemplate(w, "base", struct {
		Title string
		Body  string
	}{Title: "Dashboard", Body: "{{ .Title }}"})
	if err != nil {
		common.ErrorWriter(w, err)
	}
}

func RegisterHandlers(r *mux.Router) {
	r.HandleFunc("/dashboard", dashboard)

}
