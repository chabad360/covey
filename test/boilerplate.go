package test

import (
	"context"
	"fmt"
	"io"

	"net/http"
	"net/http/httptest"

	"github.com/go-playground/pure/v5"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest"
)

//revive:disable:function-result-limit

// Boilerplate creates a new dockertest connection.
func Boilerplate() (*dockertest.Pool, *dockertest.Resource, *pgxpool.Pool, error) {
	var err error
	var db *pgxpool.Pool
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not connect to docker: %s", err)
	}

	resource, err := pool.Run("postgres", "12", []string{"POSTGRES_PASSWORD=secret", "POSTGRES_DB=covey"})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not start resource: %s", err)
	}

	if err = pool.Retry(func() error {
		var err error
		db, err = pgxpool.Connect(context.Background(),
			fmt.Sprintf("postgres://postgres:secret@localhost:%s/%s?sslmode=disable",
				resource.GetPort("5432/tcp"), "covey"))
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, nil, nil, fmt.Errorf("could not connect to docker: %s", err)
	}

	_, err = db.Exec(context.Background(), `
	CREATE EXTENSION pgcrypto;

	CREATE TABLE nodes (
    	id TEXT PRIMARY KEY NOT NULL,
	    id_short TEXT UNIQUE NOT NULL,
    	name TEXT UNIQUE NOT NULL,
	    plugin TEXT NOT NULL,
    	details JSONB NOT NULL
	);

	CREATE TABLE tasks (
    	id TEXT PRIMARY KEY NOT NULL,
	    id_short TEXT UNIQUE NOT NULL,
    	plugin TEXT NOT NULL,
    	state INT NOT NULL,
    	node TEXT NOT NULL,
    	time TEXT,
    	exit_code INT NOT NULL,
    	log JSON,
    	details JSONB NOT NULL
	);

	CREATE TABLE jobs (
    	id TEXT PRIMARY KEY NOT NULL,
    	id_short TEXT UNIQUE NOT NULL,
    	name TEXT UNIQUE NOT NULL,
    	cron TEXT,
    	nodes JSONB NOT NULL,
    	tasks JSON NOT NULL,
    	task_history JSONB
	);

	CREATE TABLE users (
    	id SERIAL PRIMARY KEY NOT NULL,
    	username TEXT UNIQUE NOT NULL,
    	password_hash TEXT NOT NULL
	);`)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error preping the database: %s", err)
	}

	return pool, resource, db, nil
}

//revive:disable:cyclomatic

// PureBoilerplate handles boilerplate code for setting up pure.
func PureBoilerplate(method string, path string, handler http.HandlerFunc) http.Handler {
	p := pure.New()
	switch method {
	case http.MethodGet:
		p.Get(path, handler)
	case http.MethodPost:
		p.Post(path, handler)
	case http.MethodHead:
		p.Head(path, handler)
	case http.MethodPut:
		p.Put(path, handler)
	case http.MethodDelete:
		p.Delete(path, handler)
	case http.MethodConnect:
		p.Connect(path, handler)
	case http.MethodOptions:
		p.Options(path, handler)
	case http.MethodPatch:
		p.Patch(path, handler)
	case http.MethodTrace:
		p.Trace(path, handler)
	default:
		p.Handle(method, path, handler)
	}

	return p.Serve()
}

//revive:enable:cyclomatic

// HTTPBoilerplate handles boilerplate code for setting up an HTTP test.
func HTTPBoilerplate(method string, path string, body io.Reader) (*httptest.ResponseRecorder, *http.Request, error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, nil, err
	}
	rr := httptest.NewRecorder()

	return rr, req, nil
}
