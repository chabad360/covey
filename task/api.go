package task

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-playground/pure/v5"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/models"
	"github.com/chabad360/covey/storage"
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

func tasksGet(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	var q storage.Query
	err := q.Setup(r)
	common.ErrorWriter(w, err)

	var tasks interface{}

	if q.Expand {
		var t []models.Task
		err = q.Query("tasks", &t)
		tasks = t
	} else {
		var t []string
		err = q.Query("tasks", &t)
		tasks = t
	}
	common.ErrorWriter(w, err)

	common.Write(w, tasks)
}

func taskGet(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	vars := pure.RequestVars(r)

	t, ok := storage.GetTask(vars.URLParam("task"))
	common.ErrorWriter404(w, vars.URLParam("task"), ok)

	common.Write(w, t)
}

// RegisterHandlers registers the handlers for the task module.
func RegisterHandlers(r pure.IRouteGroup) {
	log.Println("Registering Task module API handlers...")

	r.Post("", taskNew)
	r.Get("", tasksGet)

	t := r.Group("/:task")
	t.Get("", taskGet)
}
