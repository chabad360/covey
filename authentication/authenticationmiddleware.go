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
		if r.URL.Path == "/auth/login" || r.URL.Path == "/login" {
			next(w, r)
			return
		}

		c, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/ui/login?url="+r.URL.Path, http.StatusFound)
			return
		}

		_, err = parseToken(c.Value, "user", "all")
		if err != nil {
			common.ErrorWriter(w, err)
			return
		}

		next(w, r)
	}
}

// AuthAPIMiddleware handles authentication for the API.
func AuthAPIMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tokenString string
		header := r.Header.Get("Authorization")
		if header != "" {
			splitToken := strings.Split(header, "Bearer ")
			tokenString = splitToken[1]
		}

		if tokenString == "" {
			common.ErrorWriterCustom(w, fmt.Errorf("forbidden"), http.StatusForbidden)
			return
		}

		claim, err := parseToken(tokenString, "api", "all")
		if err != nil { // This is here incase a user is trying to generate a token
			claim, err = parseToken(tokenString, "user", "all")
			if err != nil {
				common.ErrorWriterCustom(w, err, http.StatusUnauthorized)
				return
			}
		}

		r.Header.Add("X-User-ID", string(claim.Subject))
		next(w, r)
	}
}
