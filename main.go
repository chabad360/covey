package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/chabad360/covey/host"
)

const (
	version = "0.1"
)

func registerHandlers(r *mux.Router) {
	apiRouter := r.PathPrefix("/api/v1").Subrouter()

	apiRouter.HandleFunc("/version", getVersion)

	host.AddHandlers(apiRouter.PathPrefix("/host").Subrouter())
}
 
func main() {
	r := mux.NewRouter()
	registerHandlers(r)
	
	http.ListenAndServe(":8080", r)
}

func getVersion(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, version)
}