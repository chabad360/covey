package common

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/chabad360/covey/test"
)

func TestErrorWriter(t *testing.T) {
	rr, req, err := test.HTTPBoilerplate("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer Recover()
		ErrorWriter(w, fmt.Errorf("test"))
	}).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("ErrorWriter status = %v, want %v", status, http.StatusInternalServerError)
	}

	if rr.Body.String() != `{"error":"test"}
` {
		t.Errorf("ErrorWriter body = %v, want %v", rr.Body.String(), `{"error":"test"}`)
	}
}

func TestErrorWriter404(t *testing.T) {
	rr, req, err := test.HTTPBoilerplate("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer Recover()
		ErrorWriter404(w, "test")
	}).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("ErrorWriter404 status = %v, want %v", status, http.StatusNotFound)
	}

	if rr.Body.String() != `{"error":"404 test not found"}
` {
		t.Errorf("ErrorWriter404 body = %v, want %v", rr.Body.String(), `{"error":"404 test not found"}`)
	}
}

func TestErrorWriterCustom(t *testing.T) {
	rr, req, err := test.HTTPBoilerplate("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer Recover()
		ErrorWriterCustom(w, fmt.Errorf("test"), http.StatusUnauthorized)
	}).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("ErrorWriterCustom status = %v, want %v", status, http.StatusUnauthorized)
	}
	if rr.Body.String() != `{"error":"test"}
` {
		t.Errorf("ErrorWriterCustom body = %v, want %v", rr.Body.String(), `{"error":"test"}`)
	}
}

func TestWrite(t *testing.T) {
	rr, req, err := test.HTTPBoilerplate("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Write(w, struct{ Test string }{"test"})
	}).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Write status = %v, want %v", status, http.StatusOK)
	}

	if rr.Body.String() != `{"Test":"test"}
` {
		t.Errorf("Write body = %v, want %v", rr.Body.String(), `{"Test":"test"}
`)
	}
}
