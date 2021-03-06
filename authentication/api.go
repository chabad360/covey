package authentication

import (
	"net/http"

	"github.com/go-playground/pure/v5"

	"github.com/chabad360/covey/common"
)

func tokenGetAPI(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	id := r.Header.Get("X-User-ID")

	token, eTime, err := createToken(id, "api", []string{"all"})
	common.ErrorWriter(w, err)

	common.Write(w, struct {
		Token     string `json:"token"`
		ExpiresAt int64  `json:"expires_at"`
	}{Token: token, ExpiresAt: eTime.Unix()})
}

// RegisterAPIHandlers registers the API handlers for the authentication api.
func RegisterAPIHandlers(r pure.IRouteGroup) {
	r.Get("/token", tokenGetAPI)
}
