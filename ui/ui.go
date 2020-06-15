package ui

import (
	"fmt"
	"net/http"

	"github.com/chabad360/covey/asset"
	"github.com/chabad360/covey/common"
	"github.com/go-playground/pure/v5"
)

func dashboard(w http.ResponseWriter, _ *http.Request) {
	base := getTemplate("base")
	err := base.ExecuteTemplate(w, "base", struct {
		Title string
		Body  string
	}{Title: "Dashboard", Body: "{{ .Title }}"})
	if err != nil {
		common.ErrorWriter(w, err)
	}
}

func fsMust(f string) string {
	str, ok := asset.FS.String(f)
	if !ok {
		panic(fmt.Errorf("fsMust: invalid file %v", f))
	}
	return str
}

// RegisterHandlers registers the handlers for the ui module.
func RegisterHandlers(r pure.IRouteGroup) {
	r.Get("/dashboard", dashboard)
}
