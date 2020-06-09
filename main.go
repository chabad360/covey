package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/chabad360/covey/authentication"
	"github.com/chabad360/covey/job"
	"github.com/chabad360/covey/node"
	"github.com/chabad360/covey/storage"
	"github.com/chabad360/covey/task"
	"github.com/chabad360/covey/ui"
	"github.com/go-playground/pure/v5"
)

const (
	version = "v0.1"
)

// GetVersion returns the current version of Covey
func GetVersion(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, version)
}

// RegisterHandlers registers the core Covey API handlers
func RegisterHandlers(r pure.IRouteGroup) {
	log.Println("Registering Core module API handlers...")

	r.Get("/version", GetVersion)

	// err := r.Walk(common.Walk)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println()
}

func loadHandlers(r *pure.Mux) {
	r.Use(authentication.AuthUserMiddleware)
	r.Use(loggingMiddleware)

	ui.RegisterHandlers(r)
	authentication.RegisterHandlers(r.Group("/auth"))

	apiRouter := r.GroupWithNone("/api/v1")
	apiRouter.Use(loggingMiddleware)
	apiRouter.Use(authentication.AuthAPIMiddleware)

	RegisterHandlers(apiRouter)
	authentication.RegisterAPIHandlers(apiRouter.Group("/auth"))

	node.RegisterHandlers(apiRouter.Group("/nodes"))
	node.RegisterIndividualHandlers(apiRouter.Group("/node"))

	task.RegisterHandlers(apiRouter.Group("/tasks"))
	task.RegisterIndividualHandlers(apiRouter.Group("/task"))

	job.RegisterHandlers(apiRouter.Group("/jobs"))
	job.RegisterIndividualHandlers(apiRouter.Group("/job"))
}

func initialize() {
	storage.Init()

	// node.LoadConfig()
	// task.LoadConfig()
	job.Init()
}

func main() {
	log.Printf("Starting up Covey %s", version)

	r := pure.New()
	loadHandlers(r)

	initialize()

	fmt.Println()
	log.Println("Ready to serve!")
	fmt.Println()
	log.Fatal(http.ListenAndServe(":8080", r.Serve()))
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("API called", r.Method, r.RequestURI)
		next(w, r)
	}
}
