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
	defer common.Recover()

	reqBody, _ := ioutil.ReadAll(r.Body)
	t, err := NewTask(reqBody)
	common.ErrorWriter(w, err)

	w.Header().Add("Location", "/api/v1/task/"+t.GetIDShort())
	common.Write(w, t)
}

func taskGet(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	vars := pure.RequestVars(r)

	t, ok := getTask(vars.URLParam("task"))
	if !ok {
		common.ErrorWriter404(w, vars.URLParam("task"))
	}

	t.GetLog()

	if p := strings.Split(r.URL.Path, "/"); len(p) == 6 {
		common.Write(w, t.GetLog())
	} else {
		common.Write(w, t)
	}
}

// RegisterHandlers registers the handlers for the task module.
func RegisterHandlers(singleRoute pure.IRouteGroup, newRoute pure.IRouteGroup) {
	log.Println("Registering Task module API handlers...")

	newRoute.Post("/new", taskNew)

	t := singleRoute.Group("/:task")
	t.Get("", taskGet)
	t.Get("/log", taskGet)
}
