package authentication

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/chabad360/covey/common"
)

// AuthUserMiddleware handles authentication for users.
func AuthUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ui/auth/login" {
			next.ServeHTTP(w, r)
			return
		}

		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				common.ErrorWriterCustom(w, fmt.Errorf("Unauthorized"), http.StatusUnauthorized)
				return
			}
			common.ErrorWriter(w, err)
		}

		_, err = parseToken(c.Value)
		if err != nil {
			common.ErrorWriter(w, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// AuthAPIMiddleware handles authentication for the API.
func AuthAPIMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/auth/token" {
			next.ServeHTTP(w, r)
			return
		}

		var tokenString string
		header := r.Header.Get("Authorization")
		if header != "" {
			splitToken := strings.Split(header, "Bearer ")
			tokenString = splitToken[1]
		}

		if tokenString == "" {
			common.ErrorWriterCustom(w, fmt.Errorf("Unauthorized"), http.StatusUnauthorized)
			return
		}

		_, err := parseToken(tokenString)
		if err != nil {
			common.ErrorWriter(w, err)
			return
		}

		// if ok := claims.AllowedClaims[r.URL.Path]; !ok {
		// 	common.ErrorWriterCustom(w, fmt.Errorf("Forbidden"), http.StatusForbidden)
		// 	return
		// }

		next.ServeHTTP(w, r)
	})
}
