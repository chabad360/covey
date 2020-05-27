package task

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"plugin"
	"strings"

	"github.com/chabad360/covey/task/types"
	"github.com/gorilla/mux"
)

var (
	tasks = make(map[string]types.ITask)
)

// NewTask creates and starts a new task.
func NewTask(w http.ResponseWriter, r *http.Request) {
	var task types.Task
	reqBody, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(reqBody, &task); err != nil {
		errorWriter(w, err)
		return
	}

	p, err := loadPlugin(task.Plugin)
	if err != nil {
		errorWriter(w, err)
		return
	}

	t, err := p.NewTask(reqBody)
	if err != nil {
		errorWriter(w, err)
		return
	}
	x := t.GetID()
	tasks[hex.EncodeToString(x[:])] = t
	j, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		errorWriter(w, err)
		return
	}
	f, err := os.Create("./config/tasks.json")
	if err != nil {
		errorWriter(w, err)
		return
	}
	defer f.Close()
	if err = f.Chmod(0600); err != nil {
		errorWriter(w, err)
		return
	}
	if _, err = f.Write(j); err != nil {
		errorWriter(w, err)
		return
	}

	j, err = json.MarshalIndent(t, "", "\t")
	if err != nil {
		errorWriter(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	xs := t.GetIDShort()
	w.Header().Set("Location", "/api/v1/tasks/"+hex.EncodeToString(xs[:]))
	fmt.Fprintf(w, string(j))
}

// RegisterHandlers registers the mux handlers for the Task module.
func RegisterHandlers(r *mux.Router) {
	log.Println("Registering Task module API handlers...")
	r.HandleFunc("/new", NewTask).Methods("POST")

	err := r.Walk(walk)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()
}

func walk(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	path, err := route.GetPathTemplate()
	methods, err := route.GetMethods()
	if err == nil {
		fmt.Println("Route:", strings.Join(methods, ","), "\t", string(path))
	}
	return nil
}

func loadPlugin(pluginName string) (types.TaskPlugin, error) {
	p, err := plugin.Open("./plugins/task/" + pluginName + ".so")
	if err != nil {
		return nil, err
	}

	n, err := p.Lookup("Plugin")
	if err != nil {
		return nil, err
	}

	var s types.TaskPlugin
	s, ok := n.(types.TaskPlugin)
	if !ok {
		return nil, fmt.Errorf(pluginName, " does not provide a NodePlugin")
	}

	return s, nil
}

func errorWriter(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "{'error':'%s'}", err)
}
