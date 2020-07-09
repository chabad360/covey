package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chabad360/covey/models"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/node"
	"github.com/go-playground/pure/v5"
)

var queues = make(map[string]*TaskList)

func queueTask(nodeID string, taskID string, taskCommand string) error {
	t := AgentTask{
		ID:      taskID,
		Command: taskCommand,
	}

	id, ok := node.GetNodeIDorName(nodeID, "id")
	if !ok {
		return fmt.Errorf("%v is not a valid node ID", nodeID)
	}

	var q *TaskList
	if queues[id] == nil {
		q = &TaskList{}
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

	n, ok := node.GetNodeIDorName(vars.URLParam("node"), "id")
	if !ok {
		common.ErrorWriter404(w, vars.URLParam("node"))
	}

	b, err := ioutil.ReadAll(r.Body)
	common.ErrorWriter(w, err)

	var x TaskInfo

	err = json.Unmarshal(b, &x)
	common.ErrorWriter(w, err)

	if x.ID == "hello" {
		if n, ok = node.GetNodeIDorName(n, "name"); !ok {
			common.ErrorWriter(w, fmt.Errorf("node %s not found", n))
		}

		if err = initAgent(n); err != nil {
			common.ErrorWriter(w, err)
		}
	} else {
		saveTask(&x)
	}

	common.Write(w, queues[n])
	delete(queues, n)
}

func initQueues(tasks []models.Task) error {
	for i := range tasks {
		t := tasks[i]

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

		if err = queueTask(t.Node.ID, t.ID, cmd); err != nil {
			return err
		}
	}

	return nil
}

func initAgent(agent string) error {
	refreshDB()

	var t []models.Task
	result := db.Where("state = ?", models.StateQueued).Where("node = ?", agent).Find(&t)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

	if err := initQueues(t); err != nil {
		return err
	}

	return nil
}

// Init initializes the agent queues.
func Init() error {
	refreshDB()

	var t []models.Task
	result := db.Where("state = ?", models.StateQueued).Find(&t)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

	if err := initQueues(t); err != nil {
		return err
	}

	return nil
}

// RegisterAgentHandlers registers the handler for receiving information from agents.
func RegisterAgentHandlers(r pure.IRouteGroup) {
	r.Post("/:node", agentPost)
}
