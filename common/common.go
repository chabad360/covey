package common

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"net/http"

	json "github.com/json-iterator/go"
)

type writerError struct {
	err error
}

func (w *writerError) Error() error {
	return w.err
}

// Recover allows you to recover from a writerError and to exit a function.
func Recover() {
	if r := recover(); r != nil {
		if err, ok := r.(*writerError); ok {
			log.Println(err.Error())
			return
		}
		panic(fmt.Errorf("r: %v, type: %T", r, r))
	}
}

// RandomString generates a random 32 byte random string.
func RandomString() string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 32)

	for i := range b {
		r, err := rand.Int(rand.Reader, big.NewInt(0xffffffff))
		if err != nil {
			panic(err)
		}

		b[i] = charset[r.Int64()%int64(len(charset))]
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
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(code)
		jErr := json.NewEncoder(w).Encode(struct {
			Error string `json:"error"`
		}{err.Error()})
		if jErr != nil {
			panic(&writerError{jErr})
		}

		panic(&writerError{err})
	}
}

// Write writes the interface as a JSON to the ResponseWriter.
func Write(w http.ResponseWriter, i interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(i)
	ErrorWriter(w, err)
}

// GenerateID takes an object, converts it to text (adds two 32 byte random strings) and returns a sha256 hash from it.
func GenerateID(item interface{}) string {
	id := sha256.Sum256([]byte(fmt.Sprintf("%v%v%v", RandomString(), item, RandomString())))
	return hex.EncodeToString(id[:])
}
