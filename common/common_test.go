package common

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGenerateID(t *testing.T) {
	if got := GenerateID("simple"); got != "a7a39b72f29718e653e73503210fbb597057b7a1c77d1fe321a1afcff041d4e1" {
		t.Errorf("GenerateID() = %v, want %v", got, "a7a39b72f29718e653e73503210fbb597057b7a1c77d1fe321a1afcff041d4e1")
	}
}

func TestErrorWriter(t *testing.T) {
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ErrorWriter(w, fmt.Errorf("test"))
	})
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("ErrorWriter status = %v, want %v", status, http.StatusInternalServerError)
	}

	if rr.Body.String() != `{"error":"test"}` {
		t.Errorf("ErrorWriter body = %v, want %v", rr.Body.String(), `{"error":"test"}`)
	}
}

func TestErrorWriter404(t *testing.T) {
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ErrorWriter404(w, "test")
	})
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("ErrorWriter404 status = %v, want %v", status, http.StatusNotFound)
	}

	if rr.Body.String() != `{"error":"404 test not found"}` {
		t.Errorf("ErrorWriter404 body = %v, want %v", rr.Body.String(), `{"error":"404 test not found"}`)
	}
}

func TestErrorWriterCustom(t *testing.T) {
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ErrorWriterCustom(w, fmt.Errorf("test"), http.StatusUnauthorized)
	})
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("ErrorWriterCustom status = %v, want %v", status, http.StatusUnauthorized)
	}
	if rr.Body.String() != `{"error":"test"}` {
		t.Errorf("ErrorWriterCustom body = %v, want %v", rr.Body.String(), `{"error":"test"}`)
	}
}

func TestWrite(t *testing.T) {
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Write(w, struct {
			Test string `json:"test"`
		}{Test: "test"})
	})
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Write status = %v, want %v", status, http.StatusOK)
	}

	if rr.Body.String() != `{"test":"test"}
` {
		t.Errorf("Write body = %v, want %v", rr.Body.String(), `{"test":"test"}
`)
	}
}
