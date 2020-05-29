package task

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

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
	w.Header().Add("Content-Type", "application/json")

	var task types.Task
	reqBody, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(reqBody, &task); err != nil {
		common.ErrorWriter(w, err)
		return
	}

	p, err := loadPlugin(task.Plugin)
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}

	t, err := p.NewTask(reqBody)
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}
	tasks[t.GetID()] = t
	tasksShort[t.GetIDShort()] = t.GetID()
	j, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}
	f, err := os.Create("./config/tasks.json")
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}
	defer f.Close()
	if err = f.Chmod(0600); err != nil {
		common.ErrorWriter(w, err)
		return
	}
	if _, err = f.Write(j); err != nil {
		common.ErrorWriter(w, err)
		return
	}

	j, err = json.MarshalIndent(t, "", "\t")
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}
	w.Header().Add("Location", "/api/v1/task/"+t.GetIDShort())
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(j))
}

func taskGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	t, ok := getTask(vars["task"])
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 %v not found", vars["task"])
		return
	}

	t.GetLog()
	w.Header().Add("Content-Type", "application/json")

	j, err := json.MarshalIndent(t, "", "\t")
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(j))
}

func taskGetLog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	t, ok := getTask(vars["task"])
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 %v not found", vars["task"])
		return
	}

	w.Header().Add("Content-Type", "application/json")

	j, err := json.MarshalIndent(t.GetLog(), "", "\t")
	if err != nil {
		common.ErrorWriter(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(j))
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