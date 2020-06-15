package authentication

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chabad360/covey/common"
	"github.com/go-playground/pure/v5"
)

func tokenGet(w http.ResponseWriter, r *http.Request) {
	user := &credentials{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}
	if user.Username == "" || user.Password == "" {
		common.ErrorWriterCustom(w, fmt.Errorf("forbbiden"), http.StatusForbidden)
	}

	id, err := GetUser(*user)
	if err != nil {
		common.ErrorWriterCustom(w, fmt.Errorf("unauthorized"), http.StatusUnauthorized)
		return
	}

	token, eTime, err := createToken(id, "user", []string{"all"})
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name: "token",
		// Domain:   r.Host,
		Value:    token,
		Expires:  *eTime,
		MaxAge:   int(time.Until(*eTime).Seconds()),
		HttpOnly: true,
		Path:     "/",
	})

	if r.URL.Query().Get("url") != "" {
		http.Redirect(w, r, r.URL.Query().Get("url"), http.StatusFound)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

func tokenRemove(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		MaxAge:  -1,
		Expires: time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC),
		Path:    "/",
	})

	http.Redirect(w, r, "/login", http.StatusFound)
}

// RegisterHandlers registers the API handlers for user authentication.
func RegisterHandlers(r pure.IRouteGroup) {
	r.Post("/login", tokenGet)
	r.Get("/logout", tokenRemove)
}
