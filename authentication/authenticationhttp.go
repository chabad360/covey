package authentication

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chabad360/covey/common"
	"github.com/go-playground/pure/v5"
)

func tokenGet(w http.ResponseWriter, r *http.Request) {
	user := &credentials{}
	var ok bool
	user.Username, user.Password, ok = r.BasicAuth()
	if !ok {
		common.ErrorWriterCustom(w, fmt.Errorf("forbidden"), http.StatusForbidden)
		return
	}

	id, err := GetUser(*user)
	if err != nil {
		common.ErrorWriterCustom(w, fmt.Errorf("unauthorized"), http.StatusUnauthorized)
		return
	}

	token, eTime, err := createToken(id, "user", nil)
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
	})

	if r.URL.Query().Get("url") != "" {
		http.Redirect(w, r, r.URL.Query().Get("url"), http.StatusFound)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

func tokenRemove(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		MaxAge: -1,
	})

	http.Redirect(w, r, "/ui/login", http.StatusFound)
}

// RegisterHandlers registers the API handlers for user authentication.
func RegisterHandlers(r pure.IRouteGroup) {
	r.Get("/login", tokenGet)
	r.Post("/logout", tokenRemove)

	crashKey = randomString()
}
