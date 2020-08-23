package common

import (
	"errors"
	"fmt"
	"log"
)

// ErrWritten is returned when an error sent to an ErrorWriter has been written successfully.
var ErrWritten = errors.New("There was an error")

// ErrWriting is returned when an ErrorWriter fails to write successfully.
var ErrWriting = errors.New("error writing error")

// Recover allows you to recover from a writerError and to exit a function.
func Recover() {
	if r := recover(); r != nil {
		if err, ok := r.(error); ok && errors.Is(err, ErrWritten) {
			log.Println(err.Error())
			return
		}
		panic(fmt.Errorf("r: %v, type: %T", r, r))
	}
}
