package test

import (
	"fmt"
	"io"

	"net/http"
	"net/http/httptest"

	"github.com/chabad360/covey/models"
	"github.com/go-playground/pure/v5"
	"github.com/ory/dockertest"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//revive:disable:function-result-limit

// Boilerplate creates a new dockertest connection.
func Boilerplate() (*dockertest.Pool, *dockertest.Resource, *gorm.DB, error) {
	var err error
	var db *gorm.DB
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
		db, err = gorm.Open(postgres.Open(
			fmt.Sprintf("postgres://postgres:secret@localhost:%s/%s?sslmode=disable",
				resource.GetPort("5432/tcp"), "covey")), &gorm.Config{})
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, nil, nil, fmt.Errorf("could not connect to docker: %s", err)
	}

	db.Exec("CREATE EXTENSION pgcrypto;")

	err = db.AutoMigrate(&models.Node{}, &models.Task{}, &models.Job{}, &models.User{})
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
