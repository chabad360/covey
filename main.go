package main

import (
	"fmt"
	"net/http"

	"github.com/chabad360/covey/node"
	"github.com/gorilla/mux"
)

const (
	version = "0.1"
)

func registerHandlers(r *mux.Router) {
	apiRouter := r.PathPrefix("/api/v1").Subrouter()

	apiRouter.HandleFunc("/version", getVersion)

	node.RegisterHandlers(apiRouter.PathPrefix("/node").Subrouter())
}

func loadConfig() {
	node.LoadConfig()
}

func main() {
	r := mux.NewRouter()
	registerHandlers(r)

	loadConfig()

	http.ListenAndServe(":8080", r)
}

func getVersion(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, version)
}
