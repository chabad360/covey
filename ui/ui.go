package ui

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/chabad360/covey/asset"
	"github.com/chabad360/covey/common"
	"github.com/go-playground/pure/v5"
)

func dashboard(w http.ResponseWriter, r *http.Request) {
	base := getTemplate("dashboard")
	err := base.ExecuteTemplate(w, "base", &page{Title: "Dashboard", URL: strings.Split(r.URL.Path, "/")})
	if err != nil {
		common.ErrorWriter(w, err)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	login := getTemplate("login")
	err := login.ExecuteTemplate(w, "login", &page{Title: "Login", URL: strings.Split(r.URL.Path, "/")})
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
	r.Get("/login", login)
	r.Get("/tasks", tasks)
}
