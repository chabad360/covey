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

// ErrorWriter writes an error in the JSON format to the http.ResponseWriter.
func ErrorWriter(w http.ResponseWriter, err error) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "{'error':'%s'}", err)
}

// Walk walks the route provided to it and lists all the routes and their methods.
// Warning: This will skip over any routes that don't have a method assigned to them.
func Walk(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	path, err := route.GetPathTemplate()
	methods, err := route.GetMethods()
	if err == nil {
		fmt.Println("Route:", strings.Join(methods, ","), "\t", string(path))
	}
	return nil
}

// GenerateID takes a object and converts it to json and then returns a sha256 hash of the object.
// Warning: This is not guaranteed to be unique, please ensure that your object includes a field that is unique (i.e. time.Now).
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
