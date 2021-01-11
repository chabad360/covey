package ui

import (
	"fmt"
	"net/http"

	"github.com/chabad360/covey/common"
)

// ErrorWriter writes an UI error with the 500 status code.
func ErrorWriter(w http.ResponseWriter, err error) {
	ErrorWriterCustom(w, err, http.StatusInternalServerError, "Uh Oh! There was an error!")
}

// ErrorWriter404 writes an UI error with a 404 code.
func ErrorWriter404(w http.ResponseWriter, name string, ok bool) {
	if !ok {
		ErrorWriterCustom(w, fmt.Errorf("%v not found", name), http.StatusNotFound, "What you're looking for doesn't exist...")
	}
}

// ErrorWriterCustom writes an UI error to the http.ResponseWriter with a custom status code.
func ErrorWriterCustom(w http.ResponseWriter, err error, code int, title string) {
	if err == nil {
		return
	}

	w.WriteHeader(code)

	p := &Page{
		Title: title,
		Details: struct {
			Code int
			Text string
		}{code, err.Error()},
	}

	t := GetTemplate("ec")
	wErr := t.ExecuteTemplate(w, "ec", p)
	if wErr != nil {
		panic(fmt.Errorf("%w: %v", common.ErrWriting, wErr))
	}

	panic(fmt.Errorf("%w: %v", common.ErrWritten, err))
}
