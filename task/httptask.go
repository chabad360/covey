package task

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/task/types"
	"github.com/gorilla/mux"
)

var (
	tasks      = make(map[string]types.ITask)
	tasksShort = make(map[string]string)
)

// TaskNew creates and starts a new task.
func taskNew(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	t, err := NewTask(reqBody)
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}

	w.Header().Add("Location", "/api/v1/task/"+t.GetIDShort())
	common.Write(w, t)
}

func taskGet(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	t, ok := GetTask(vars["task"])
	if !ok {
		common.ErrorWriter404(w, vars["task"])
		return
	}

	t.GetLog()
	SaveTask(t)
	common.Write(w, t)
}

func taskGetLog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	t, ok := GetTask(vars["task"])
	if !ok {
		common.ErrorWriter404(w, vars["task"])
		return
	}

	t.GetLog()
	SaveTask(t)
	common.Write(w, t.GetLog())
}

// RegisterHandlers registers the mux handlers for the Task module.
func RegisterHandlers(r *mux.Router) {
	log.Println("Registering Task module API handlers...")

	r.HandleFunc("/new", taskNew).Methods("POST")
	r.HandleFunc("/{task}", taskGet).Methods("GET")
	r.HandleFunc("/{task}/log", taskGetLog).Methods("GET")

	err := r.Walk(common.Walk)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()
}
