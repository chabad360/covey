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
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "/ui/login?url="+r.URL.Path, http.StatusFound)
				return
			}
			common.ErrorWriter(w, err)
		}

		_, err = parseToken(c.Value)
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
			common.ErrorWriterCustom(w, fmt.Errorf("unauthorized"), http.StatusUnauthorized)
			return
		}

		claim, err := parseToken(tokenString)
		if err != nil {
			common.ErrorWriter(w, err)
			return
		}

		// if ok := claims.AllowedClaims[r.URL.Path]; !ok {
		// 	common.ErrorWriterCustom(w, fmt.Errorf("Forbidden"), http.StatusForbidden)
		// 	return
		// }

		r.Header.Add("X-User-ID", string(claim.UserID))
		next(w, r)
	}
}
