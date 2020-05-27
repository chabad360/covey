package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/chabad360/covey/node"
	"github.com/gorilla/mux"
)

const (
	version = "0.1"
)

// GetVersion returns the current version of Covey
func GetVersion(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, version)
}

// RegisterHandlers registers the core Covey API handlers
func RegisterHandlers(r *mux.Router) {
	r.HandleFunc("/version", GetVersion)
}

func registerHandlers(r *mux.Router) {
	apiRouter := r.PathPrefix("/api/v1").Subrouter()

	registerHandlers(apiRouter)
	node.RegisterHandlers(apiRouter.PathPrefix("/node").Subrouter())
}

func loadConfig() {
	node.LoadConfig()
}

func main() {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	registerHandlers(r)

	err := r.Walk(walk)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println()

	loadConfig()

	http.ListenAndServe(":8080", r)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("API called", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func walk(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	path, err := route.GetPathTemplate()
	methods, err := route.GetMethods()
	if err == nil {
		fmt.Println("Registered route:", strings.Join(methods, ","), "\t", string(path))
	}
	return nil
}
