package task

import (
	"encoding/json"
	"fmt"
	"github.com/chabad360/covey/models"
	"github.com/chabad360/covey/storage"
	"github.com/go-playground/pure/v5"
	"io/ioutil"
	"net/http"

	"github.com/chabad360/covey/common"
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
		n, ok = storage.GetNodeIDorName(n, "name")
		common.ErrorWriter404(w, n, ok)

		common.ErrorWriter(w, Init(n))
	} else {
		common.ErrorWriter(w, storage.SaveTask(&x))
	}

	common.Write(w, queues[n])
	delete(queues, n)
}

func initQueues(tasks []models.Task) error {
	for _, t := range tasks {
		p, err := loadPlugin(t.Plugin)
		if err != nil {
			return err
		}

		cmd, err := p.GetCommand(t)
		if err != nil {
			return err
		}

		if err = queueTask(t.Node, t.ID, cmd); err != nil {
			return err
		}
	}

	return nil
}

// Init initializes the agent queues.
func Init(agent string) error {
	var t []models.Task
	tx := storage.DB.Where("state = ?", models.StateQueued)

	if agent != "" {
		tx.Where("node = ?", agent)
	}

	if err := tx.Find(&t).Error; err != nil {
		return err
	}

	return initQueues(t)
}

// RegisterAgentHandlers registers the handler for receiving information from agents.
func RegisterAgentHandlers(r pure.IRouteGroup) {
	r.Post("/:node", agentPost)
}
