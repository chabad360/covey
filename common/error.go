package common

import (
	"fmt"
	"log"
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
