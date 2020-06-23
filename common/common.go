package common

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	json "github.com/json-iterator/go"
)

var (
	random = rand.New(
		rand.NewSource(time.Now().UnixNano()))
)

func RandomString() string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 32)
	for i := range b {
		b[i] = charset[random.Int63()%int64(len(charset))]
	}
	return string(b)
}

// ErrorWriter writes an error in the JSON format to the http.ResponseWriter.
func ErrorWriter(w http.ResponseWriter, err error) {
	ErrorWriterCustom(w, err, http.StatusInternalServerError)
}

// ErrorWriter404 writes an error in the JSON format to the with a 404 code.
func ErrorWriter404(w http.ResponseWriter, name string) {
	ErrorWriterCustom(w, fmt.Errorf("404 %v not found", name), http.StatusNotFound)
}

// ErrorWriterCustom writes an error in the JSON format to the http.ResponseWriter with a custom status code.
func ErrorWriterCustom(w http.ResponseWriter, err error, code int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(struct {
		Error string `json:"error"`
	}{err.Error()})
}

// Write writes the interface as a JSON to the ResponseWriter.
func Write(w http.ResponseWriter, i interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(i)
}

// GenerateID takes a object and converts it to text and then returns a sha256 hash of the object.
// Warning: This is not guaranteed to be unique, so include a field that is unique (i.e. time.Now).
func GenerateID(item interface{}) string {
	id := sha256.Sum256([]byte(fmt.Sprintf("%v", item)))
	return hex.EncodeToString(id[:])
}
