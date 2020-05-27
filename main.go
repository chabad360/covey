package main

import (
	"fmt"
	"log"
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
	r.Use(loggingMiddleware)

	loadConfig()

	http.ListenAndServe(":8080", r)
}

func getVersion(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, version)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("API called", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
