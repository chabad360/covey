package authentication

import (
	"fmt"
	"github.com/chabad360/covey/common"
	"net/http"
	"strings"
)

// AuthUserMiddleware handles authentication for users.
func AuthUserMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/logout" || r.URL.Path == "/login" {
			next(w, r)
			return
		}

		c, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/login?url="+r.URL.Path, http.StatusFound)
			return
		}

		_, err = parseToken(c.Value, "user", "all")
		if err != nil {
			// if it's a bad cookie, we log them out, effectively deleting the cookie.
			http.Redirect(w, r, "/logout", http.StatusFound)
			return
		}

		next(w, r)
	}
}

// AuthAPIMiddleware handles authentication for the API.
func AuthAPIMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer common.Recover()

		splitToken := strings.Split(r.Header.Get("Authorization"), "Bearer")
		if len(splitToken) != 2 {
			common.ErrorWriterCustom(w, fmt.Errorf("invalid bearer token"), http.StatusBadRequest)
		}

		tokenString := strings.TrimSpace(splitToken[1])

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
