package authentication

import (
	"net/http"
	"strconv"

	"github.com/chabad360/covey/common"
	"github.com/gorilla/mux"
)

func tokenGetAPI(w http.ResponseWriter, r *http.Request) {
	ids := r.Header.Get("X-User-ID")
	id, err := strconv.Atoi(ids)
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}

	token, eTime, err := createToken(uint32(id), "api", nil)
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
func RegisterAPIHandlers(r *mux.Router) {
	r.HandleFunc("/token", tokenGetAPI).Methods("GET")
}
