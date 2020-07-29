package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chabad360/covey/models"
	"github.com/chabad360/covey/storage"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"

	"github.com/chabad360/covey/common"
	"github.com/go-playground/pure/v5"
)

var queues = make(map[string]*List)

func queueTask(nodeID string, taskID string, taskCommand string) error {
	t := agentTask{
		ID:      taskID,
		Command: taskCommand,
	}

	id, ok := storage.GetNodeIDorName(nodeID, "id")
	if !ok {
		return fmt.Errorf("%v is not a valid node ID", nodeID)
	}

	var q *List
	if queues[id] == nil {
		q = &List{}
	} else {
		q = queues[id]
	}

	q.PushBack(t)

	queues[id] = q

	return nil
}

func agentPost(w http.ResponseWriter, r *http.Request) {
	defer common.Recover()

	vars := pure.RequestVars(r)

	n, ok := storage.GetNodeIDorName(vars.URLParam("node"), "id")
	common.ErrorWriter404(w, vars.URLParam("node"), ok) // TODO: disable the agent if there is no such node

	b, err := ioutil.ReadAll(r.Body)
	common.ErrorWriter(w, err)

	var x storage.TaskInfo

	err = json.Unmarshal(b, &x)
	common.ErrorWriter(w, err)

	if x.ID == "hello" {
		if n, ok = storage.GetNodeIDorName(n, "name"); !ok {
			common.ErrorWriter(w, fmt.Errorf("node %s not found", n))
		}

		if err = initAgent(n); err != nil {
			common.ErrorWriter(w, err)
		}
	} else {
		err = storage.SaveTask(&x)
		common.ErrorWriter(w, err)
	}

	common.Write(w, queues[n])
	delete(queues, n)
}

func initQueues(tasks []models.Task) error {
	for _, t := range tasks {
		j, err := json.Marshal(t)
		if err != nil {
			return err
		}

		p, err := loadPlugin(t.Plugin)
		if err != nil {
			return err
		}

		cmd, err := p.GetCommand(j)
		if err != nil {
			return err
		}

		if err = queueTask(t.Node, t.ID, cmd); err != nil {
			return err
		}
	}

	return nil
}

func initAgent(agent string) error {
	var t []models.Task
	result := storage.DB.Where("state = ?", models.StateQueued).Where("node = ?", agent).Find(&t)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

	return initQueues(t)
}

// Init initializes the agent queues.
func Init() error {
	var t []models.Task
	result := storage.DB.Where("state = ?", models.StateQueued).Find(&t)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

	return initQueues(t)
}

// RegisterAgentHandlers registers the handler for receiving information from agents.
func RegisterAgentHandlers(r pure.IRouteGroup) {
	r.Post("/:node", agentPost)
}
