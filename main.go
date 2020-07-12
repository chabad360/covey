package main

// Make sure to run resources -declare -package=asset -output=asset/asset.go -tag="\!live" -trim assets/ assets/*

import (
	"fmt"
	"github.com/chabad360/covey/models"
	"gorm.io/gorm"
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

// GetVersion returns the current version of Covey.
func getVersion(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "%s", version)
}

func registerHandlers(r pure.IRouteGroup) {
	log.Println("Registering Core module API handlers...")

	r.Get("/version", getVersion)
}

func loadHandlers(r *pure.Mux) {
	r.Use(loggingMiddleware)
	r.Use(authentication.AuthUserMiddleware)
	r.Use(options)

	r.Get("/src/*", func() http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "max-age=2592000")
			http.FileServer(asset.FS).ServeHTTP(w, r)
		}
	}()) // Make static files cached

	// TODO: Clean up
	r.RegisterAutomaticOPTIONS(options)
	ui.RegisterHandlers(r)
	authentication.RegisterUIHandlers(r)

	job.RegisterUIHandlers(r.Group("/jobs"))
	r.Get("/new/job", job.UIJobNew) // BAD

	node.RegisterUIHandlers(r.Group("/nodes"))
	r.Get("/new/node", node.UINodeNew) // BAD

	task.RegisterUIHandlers(r.Group("/tasks"))
	r.Get("/new/task", task.UITaskNew) // BAD

	agent := r.GroupWithNone("/agent")
	agent.Use(loggingMiddleware)
	task.RegisterAgentHandlers(agent)

	apiRouter := r.GroupWithNone("/api/v1")
	apiRouter.Use(loggingMiddleware)
	apiRouter.Use(options)
	apiRouter.Use(authentication.AuthAPIMiddleware)

	registerHandlers(apiRouter)
	authentication.RegisterAPIHandlers(apiRouter.Group("/auth"))

	node.RegisterHandlers(apiRouter.Group("/nodes"))
	task.RegisterHandlers(apiRouter.Group("/tasks"))
	job.RegisterHandlers(apiRouter.Group("/jobs"))
}

func initialize() {
	job.Init()

	if err := task.Init(); err != nil {
		log.Fatal(err)
	}

	// Ensure files are available
	if asset.FS == nil {
		log.Fatal(`Remember to run 
		'resources -declare -package=asset -output=asset/asset.go -tag="\!live" -trim assets/ assets/*'`)
	}

	if _, err := asset.FS.Open("/base/base.html"); err != nil {
		log.Fatalf("Failed to open filesystem: %v", err)
	}
}

func initStorage() (*gorm.DB, error) {
	log.Println("Setup DB")

	err := storage.Init()
	if err != nil {
		return nil, err
	}
	db := storage.DB

	db.Exec("CREATE EXTENSION pgcrypto;")

	err = db.AutoMigrate(&models.Node{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Task{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Job{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	log.Printf("Starting up Covey %s", version)
	fmt.Println()

	r := pure.New()
	loadHandlers(r)

	_, err := initStorage()
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()

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

func options(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		o := r.Header.Get("origin")
		w.Header().Set("Access-Control-Allow-Origin", o)
		w.Header().Set("Access-Control-Allow-Headers", "*")
		next(w, r)
	}
}
