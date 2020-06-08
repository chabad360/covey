package authentication

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chabad360/covey/common"
	"github.com/gorilla/mux"
)

func tokenGetAPI(w http.ResponseWriter, r *http.Request) {
	user := &credentials{}
	var ok bool
	user.Username, user.Password, ok = r.BasicAuth()
	if !ok {
		common.ErrorWriterCustom(w, fmt.Errorf("Unauthorized"), http.StatusUnauthorized)
		return
	}

	id, err := GetUser(*user)
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}

	token, eTime, err := createToken(id, "api", nil)
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Token     string `json:"token"`
		ExpiresAt int64  `json:"expires_at"`
	}{Token: token, ExpiresAt: eTime.Unix()})
}

// RegisterAPIHandlers registers the API handlers for the authentication api.
func RegisterAPIHandlers(r *mux.Router) {
	r.HandleFunc("/token", tokenGetAPI).Methods("GET")

	crashKey = randomString()
}
