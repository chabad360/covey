package authentication

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/chabad360/covey/common"
)

// AuthUserMiddleware handles authentication for users.
func AuthUserMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			// r.Cookie only throws http.ErrNoCookie to we need to let someone login, only if there is no cookie.
			if r.URL.Path == "/login" {
				next(w, r)
				return
			}

			http.Redirect(w, r, "/login?url="+r.URL.Path, http.StatusFound)
			return
		}

		_, err = parseToken(c.Value, "user", "all")
		if err != nil {
			// if it's a bad cookie, we log them out, effectively deleting the cookie.
			if r.URL.Path == "/logout" {
				next(w, r)
				return
			}

			http.Redirect(w, r, "/logout", http.StatusFound)
			return
		}

		if r.URL.Path == "/login" {
			http.Redirect(w, r, "/dashboard", http.StatusFound)
			return
		}

		next(w, r)
	}
}

// AuthAPIMiddleware handles authentication for the API.
func AuthAPIMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer common.Recover()

		var tokenString string
		header := r.Header.Get("Authorization")
		if header != "" {
			splitToken := strings.Split(header, "Bearer ")
			tokenString = splitToken[1]
		}

		if tokenString == "" {
			common.ErrorWriterCustom(w, fmt.Errorf("forbidden"), http.StatusForbidden)
		}

		claim, err := parseToken(tokenString, "api", "all")
		if err != nil {
			// This is here in case a user is trying to generate a token or use the api directly.
			claim, err = parseToken(tokenString, "user", "all")
			if err != nil {
				common.ErrorWriterCustom(w, err, http.StatusUnauthorized)
			}
		}

		r.Header.Add("X-User-ID", claim.Subject)
		next(w, r)
	}
}
