package task

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"plugin"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/task/types"
	"github.com/gorilla/mux"
)

var (
	tasks     = make(map[string]types.ITask)
	taskShort = make(map[string]string)
)

// NewTask creates and starts a new task.
func NewTask(w http.ResponseWriter, r *http.Request) {
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
	taskShort[t.GetIDShort()] = t.GetID()
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
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Location", "/api/v1/tasks/"+t.GetIDShort())
	fmt.Fprintf(w, string(j))
}

// RegisterHandlers registers the mux handlers for the Task module.
func RegisterHandlers(r *mux.Router) {
	log.Println("Registering Task module API handlers...")
	r.HandleFunc("/new", NewTask).Methods("POST")

	err := r.Walk(common.Walk)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()
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
		return nil, fmt.Errorf(pluginName, " does not provide a TaskPlugin")
	}

	return s, nil
}
