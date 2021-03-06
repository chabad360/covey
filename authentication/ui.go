package authentication

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/pure/v5"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/models"
	"github.com/chabad360/covey/ui"
)

func login(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	l := ui.GetTemplate("login")
	err := l.ExecuteTemplate(w, "login",
		&ui.Page{Title: "Login", URL: strings.Split(r.URL.Path, "/"), Details: true})
	common.ErrorWriter(w, err)
}

func loginPost(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	loginPage := ui.GetTemplate("login")

	cookie, err := tokenCookie(&models.User{Username: r.FormValue("username"), Password: r.FormValue("password")})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		err = loginPage.ExecuteTemplate(w, "login",
			&ui.Page{Title: "Login", URL: strings.Split(r.URL.Path, "/"), Details: false})
		common.ErrorWriter(w, err)
	}

	http.SetCookie(w, cookie)

	if r.URL.Query().Get("url") != "" {
		http.Redirect(w, r, r.URL.Query().Get("url"), http.StatusFound)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

func logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		MaxAge:  -1,
		Expires: time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC),
		Path:    "/",
	})

	http.Redirect(w, r, "/login", http.StatusFound)
}

// RegisterUIHandlers registers the UI handlers for user authentication.
func RegisterUIHandlers(r pure.IRouteGroup) {
	r.Post("/login", loginPost)
	r.Get("/login", login)
	r.Get("/logout", logout)
}
