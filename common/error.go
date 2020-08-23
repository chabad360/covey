package common

import (
	"errors"
	"fmt"
	"log"
)

var WriterError = errors.New("There was an error")

// Recover allows you to recover from a writerError and to exit a function.
func Recover() {
	if r := recover(); r != nil {
		if err, ok := r.(error); ok && errors.Is(err, WriterError) {
			log.Println(err.Error())
			return
		}
		panic(fmt.Errorf("r: %v, type: %T", r, r))
	}
}
