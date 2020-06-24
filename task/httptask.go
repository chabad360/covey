package task

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/chabad360/covey/common"
	"github.com/go-playground/pure/v5"
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
	vars := pure.RequestVars(r)
	t, ok := GetTask(vars.URLParam("task"))
	if !ok {
		common.ErrorWriter404(w, vars.URLParam("task"))
		return
	}
	t.GetLog()

	if p := strings.Split(r.URL.Path, "/"); len(p) == 6 {
		common.Write(w, t.GetLog())
	} else {
		common.Write(w, t)
	}
}

// RegisterHandlers registers the handlers for the task module.
func RegisterHandlers(r pure.IRouteGroup) {
	log.Println("Registering Task module API handlers...")

	r.Post("/new", taskNew)
}

// RegisterIndividualHandlers registers the mux handlers for the Task module.
func RegisterIndividualHandlers(r pure.IRouteGroup) {
	t := r.Group("/:task")
	t.Get("", taskGet)
	t.Get("/log", taskGet)
}
