package common

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// ErrorWriter writes an error to the http.ResponseWriter
func ErrorWriter(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "{'error':'%s'}", err)
}

// Walk walks the route
func Walk(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	path, err := route.GetPathTemplate()
	methods, err := route.GetMethods()
	if err == nil {
		fmt.Println("Route:", strings.Join(methods, ","), "\t", string(path))
	}
	return nil
}

func GenerateID(item interface{}) (*string, error) {
	j, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	ida := sha256.Sum256(j)
	id := ida[:]
	ids := hex.EncodeToString(id)
	return &ids, nil
}
