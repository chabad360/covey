package task

import (
	"fmt"
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

func tasksGet(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()
	refreshDB()

	log.Println(r.URL.Query())
	var q common.QueryParams
	err := pure.DecodeQueryParams(r, 1, &q)
	common.ErrorWriter(w, err)

	var tasks []string
	result := db.Select("id").Find(&tasks)
	common.ErrorWriter(w, result.Error)

	common.Write(w, tasks)
}

func taskGet(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	vars := pure.RequestVars(r)

	t, err := getTask(vars.URLParam("task"))
	common.ErrorWriter404(w, fmt.Sprintf("%v", err))

	if p := strings.Split(r.URL.Path, "/"); len(p) == 6 {
		common.Write(w, t.Log)
	} else {
		common.Write(w, t)
	}
}

// RegisterHandlers registers the handlers for the task module.
func RegisterHandlers(r pure.IRouteGroup) {
	log.Println("Registering Task module API handlers...")

	r.Post("", taskNew)
	r.Get("", tasksGet)
	t := r.Group("/:task")
	t.Get("", taskGet)
	t.Get("/log", taskGet)
}
