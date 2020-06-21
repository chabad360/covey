package main

// Make sure to run resources -declare -package=asset -output=asset/asset.go -tag="\!live" -trim assets/ assets/*

import (
	"fmt"
	"log"
	"net/http"

	"github.com/chabad360/covey/asset"
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
func getVersion(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, version)
}

// RegisterHandlers registers the core Covey API handlers
func RegisterHandlers(r pure.IRouteGroup) {
	log.Println("Registering Core module API handlers...")

	r.Get("/version", getVersion)
}

func loadHandlers(r *pure.Mux) {
	r.Use(loggingMiddleware)
	r.Use(authentication.AuthUserMiddleware)
	r.Get("/src/*", http.FileServer(asset.FS).ServeHTTP)

	ui.RegisterHandlers(r)
	authentication.RegisterUIHandlers(r)
	task.RegisterUIHandlers(r.Group("/tasks"))
	r.Get("/new/task", task.UITaskNew)
	job.RegisterUIHandlers(r.Group("/jobs"))
	r.Get("/new/job", job.UIJobNew)

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
	job.Init()

	// Ensure files are available
	if asset.FS == nil {
		log.Fatal(`Remember to run 
		'resources -declare -package=asset -output=asset/asset.go -tag="\!live" -trim assets/ assets/*'`)
	}
	if _, err := asset.FS.Open("/base/base.html"); err != nil {
		log.Fatalf("Failed to open filesystem: %v", err)
	}
}

func main() {
	log.Printf("Starting up Covey %s", version)
	fmt.Println()

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
		log.Println("Called", r.Method, r.RequestURI)
		next(w, r)
	}
}
