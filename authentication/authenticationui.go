package authentication

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/ui"
	"github.com/go-playground/pure/v5"
)

func login(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	l := ui.GetTemplate("login")
	err := l.ExecuteTemplate(w, "login",
		&ui.Page{Title: "Login", URL: strings.Split(r.URL.Path, "/"), Details: true})
	common.ErrorWriter(w, err)
}

func loginP(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	login := ui.GetTemplate("login")

	cookie, err := tokenGet(&credentials{Username: r.FormValue("username"), Password: r.FormValue("password")})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		err = login.ExecuteTemplate(w, "login",
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

func tokenGet(user *credentials) (*http.Cookie, error) {
	id, err := GetUser(*user)
	if err != nil {
		return nil, err
	}

	token, eTime, err := createToken(id, "user", []string{"all"})
	if err != nil {
		return nil, fmt.Errorf("internal error")
	}

	return &http.Cookie{
		Name: "token",
		// Domain:   r.Host,
		Value:   token,
		Expires: *eTime,
		MaxAge:  int(time.Until(*eTime).Seconds()),
		Path:    "/",
	}, nil
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
	r.Post("/login", loginP)
	r.Get("/login", login)
	r.Get("/logout", logout)
}
