package authentication

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gbrlsnchs/jwt/v3"

	"github.com/chabad360/covey/common"
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

		j, err := parseToken(c.Value, "user", "all")
		if err != nil {
			log.Println(j, err)
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
		var claim *jwt.Payload
		var err error

		splitToken := strings.Split(r.Header.Get("Authorization"), "Bearer")
		if len(splitToken) != 2 {
			common.ErrorWriterCustom(w, fmt.Errorf("invalid bearer token"), http.StatusBadRequest)
		}

		tokenString := strings.TrimSpace(splitToken[1])

		claim, err = parseToken(tokenString, "api", "all")
		if err != nil {
			// This is here in case a user is trying to generate a token or use the api directly.
			claim, err = parseToken(tokenString, "user", "all")
			if err != nil {
				common.ErrorWriterCustom(w, err, http.StatusUnauthorized)
			}
		}

		log.Println(claim.Subject)
		r.Header.Add("X-User-ID", claim.Subject)
		next(w, r)
	}
}
