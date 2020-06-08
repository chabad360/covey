package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/chabad360/covey/authentication"
	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/job"
	"github.com/chabad360/covey/node"
	"github.com/chabad360/covey/storage"
	"github.com/chabad360/covey/task"
	"github.com/chabad360/covey/ui"
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

	err := r.Walk(common.Walk)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()
}

func loadHandlers(r *mux.Router) {
	apiRouter := r.PathPrefix("/api/v1").Subrouter()
	apiRouter.Use(authentication.AuthAPIMiddleware)

	RegisterHandlers(apiRouter)
	authentication.RegisterAPIHandlers(apiRouter.PathPrefix("/auth").Subrouter())
	node.RegisterHandlers(apiRouter.PathPrefix("/node").Subrouter())
	task.RegisterHandlers(apiRouter.PathPrefix("/task").Subrouter())
	job.RegisterHandlers(apiRouter.PathPrefix("/job").Subrouter())

	uiRouter := r.PathPrefix("/ui").Subrouter()
	uiRouter.Use(authentication.AuthUserMiddleware)
	ui.RegisterHandlers(uiRouter)
	authentication.RegisterHandlers(uiRouter.PathPrefix("/auth").Subrouter())
}

func loadConfig() {
	// node.LoadConfig()
	// task.LoadConfig()
	job.Init()
}

func main() {
	storage.Init()

	log.Printf("Starting up Covey %s", version)
	r := mux.NewRouter()
	loadHandlers(r)

	r.Use(loggingMiddleware)

	loadConfig()

	fmt.Println()
	log.Println("Ready to serve!")
	fmt.Println()
	log.Fatal(http.ListenAndServe(":8080", r))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("API called", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
