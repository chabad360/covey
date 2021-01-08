package task

import (
	"fmt"
	"github.com/chabad360/covey/models"
	"github.com/chabad360/covey/plugin"
	"github.com/chabad360/covey/storage"
	"github.com/go-playground/pure/v5"
	json "github.com/json-iterator/go"
	"io/ioutil"
	"log"
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
	n, ok := storage.GetNode(vars.URLParam("node"))
	//common.ErrorWriter404(w, vars.URLParam("node"), ok) // TODO: disable the agent if there is no such node
	if !ok {
		log.Printf("node %v not found", vars.URLParam("node"))
		common.Write(w, nil)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	//common.ErrorWriter(w, err)
	log.Print(err)

	var x storage.TaskInfo
	log.Print(json.Unmarshal(b, &x))

	if x.ID == "hello" {
		log.Print(Init(n))
	} else if x.ID != "" {
		log.Print(storage.SaveTask(&x))
	}

	common.Write(w, queues[n.ID])
	delete(queues, n.ID)
}

func initQueues(tasks []models.Task) error {
	for _, t := range tasks {
		p, err := plugin.GetTaskPlugin(t.Plugin)
		if err != nil {
			return err
		}

		cmd, err := p.GetCommand(t.ToSafe())
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
func Init(agent *models.Node) error {
	var t []models.Task
	tx := storage.DB.Where("state = ?", models.StateQueued)

	if agent != nil {
		tx.Where("node = ?", agent.Name)
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
