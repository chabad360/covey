package authentication

import (
	"net/http"

	"github.com/chabad360/covey/common"
	"github.com/go-playground/pure/v5"
)

func tokenGetAPI(w http.ResponseWriter, r *http.Request) {
	id := r.Header.Get("X-User-ID")

	token, eTime, err := createToken(id, "api", []string{"all"})
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}

	common.Write(w, struct {
		Token     string `json:"token"`
		ExpiresAt int64  `json:"expires_at"`
	}{Token: token, ExpiresAt: eTime.Unix()})
}

// RegisterAPIHandlers registers the API handlers for the authentication api.
func RegisterAPIHandlers(r pure.IRouteGroup) {
	r.Get("/token", tokenGetAPI)
}
