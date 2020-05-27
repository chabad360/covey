package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/chabad360/covey/node"
	"github.com/chabad360/covey/task"
	"github.com/gorilla/mux"
)

const (
	version = "v0.1"
)

// GetVersion returns the current version of Covey
func GetVersion(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, version)
}

// RegisterHandlers registers the core Covey API handlers
func RegisterHandlers(r *mux.Router) {
	log.Println("Registering Core module API handlers...")

	r.HandleFunc("/version", GetVersion).Methods("GET")

	err := r.Walk(walk)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()
}

func loadHandlers(r *mux.Router) {
	apiRouter := r.PathPrefix("/api/v1").Subrouter()

	RegisterHandlers(apiRouter)
	node.RegisterHandlers(apiRouter.PathPrefix("/node").Subrouter())
	task.RegisterHandlers(apiRouter.PathPrefix("/task").Subrouter())
}

func loadConfig() {
	node.LoadConfig()
}

func main() {
	log.Printf("Starting up Covey %s", version)
	r := mux.NewRouter()
	loadHandlers(r)

	r.Use(loggingMiddleware)

	loadConfig()

	fmt.Println()
	log.Println("Ready to serve!")
	log.Fatal(http.ListenAndServe(":8080", r))
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
		fmt.Println("Route:", strings.Join(methods, ","), "\t", string(path))
	}
	return nil
}
