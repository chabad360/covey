package ui

import (
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/go-playground/pure/v5"

	"github.com/chabad360/covey/assets"
	"github.com/chabad360/covey/common"
)

func dashboard(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	base := GetTemplate("dashboard")
	err := base.ExecuteTemplate(w, "base", &Page{Title: "Dashboard", URL: strings.Split(r.URL.Path, "/")})
	ErrorWriter(w, err)
}

func fsMust(f string) string {
	byt, err := fs.ReadFile(assets.Content, f)
	if err != nil {
		panic(fmt.Errorf("fsMust: %w", err))
	}
	return string(byt)
}

// RegisterHandlers registers the handlers for the ui module.
func RegisterHandlers(r pure.IRouteGroup) {
	r.Get("/dashboard", dashboard)
	r.Get("/", dashboard)
}
