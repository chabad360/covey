package common

import (
	"fmt"
	"net/http"

	json "github.com/json-iterator/go"
)

// ErrorWriter writes an error in the JSON format to the http.ResponseWriter.
func ErrorWriter(w http.ResponseWriter, err error) {
	ErrorWriterCustom(w, err, http.StatusInternalServerError)
}

// ErrorWriter404 writes an error in the JSON format to the with a 404 code.
func ErrorWriter404(w http.ResponseWriter, name string, ok bool) {
	if !ok {
		ErrorWriterCustom(w, fmt.Errorf("404 %v not found", name), http.StatusNotFound)
	}
}

// ErrorWriterCustom writes an error in the JSON format to the http.ResponseWriter with a custom status code.
func ErrorWriterCustom(w http.ResponseWriter, err error, code int) {
	if err == nil {
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	jErr := json.NewEncoder(w).Encode(struct {
		Error string `json:"error"`
	}{err.Error()})
	if jErr != nil {
		panic(fmt.Errorf("error writing response: %w", jErr))
	}

	panic(fmt.Errorf("%w: %v", WriterError, err))
}

// Write writes the interface as a JSON to the ResponseWriter.
func Write(w http.ResponseWriter, i interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(i)
	ErrorWriter(w, err)
}
