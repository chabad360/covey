package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/storage"
	"github.com/chabad360/covey/task/types"
	"github.com/gorilla/mux"
)

var (
	tasks      = make(map[string]types.ITask)
	tasksShort = make(map[string]string)
)

// TaskNew creates and starts a new task.
func taskNew(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	reqBody, _ := ioutil.ReadAll(r.Body)
	t, err := NewTask(reqBody)
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}

	tasks[t.GetID()] = t
	tasksShort[t.GetIDShort()] = t.GetID()
	SaveTask(t)

	z, err := storage.GetItem("tasks", t.GetID(), t)
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}

	j, err := json.MarshalIndent(z, "", "\t")
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}
	w.Header().Add("Location", "/api/v1/task/"+t.GetIDShort())
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(j))
}

func taskGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	vars := mux.Vars(r)
	t, ok := GetTask(vars["task"])
	if !ok {
		common.ErrorWriter(w, errors.New("404 not found"))
		return
	}

	t.GetLog()
	SaveTask(t)

	j, err := json.MarshalIndent(t, "", "\t")
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(j))
}

func taskGetLog(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	vars := mux.Vars(r)
	t, ok := GetTask(vars["task"])
	if !ok {
		common.ErrorWriter(w, errors.New("404 not found"))
		return
	}

	t.GetLog()
	SaveTask(t)

	j, err := json.MarshalIndent(t.GetLog(), "", "\t")
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(j))
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
