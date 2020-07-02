package task

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/node"
	"github.com/chabad360/covey/task/types"
	"github.com/go-playground/pure/v5"
)

var queues = make(map[string]*types.TaskList)

func queueTask(nodeID string, taskID string, taskCommand string) error {
	t := types.AgentTask{
		ID:      taskID,
		Command: taskCommand,
	}

	id, ok := node.GetNodeID(nodeID)
	if !ok {
		return fmt.Errorf("%v is not a valid node ID", nodeID)
	}

	var q *types.TaskList
	if queues[id] == nil {
		q = &types.TaskList{}
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

	n, ok := node.GetNodeID(vars.URLParam("node"))
	if !ok {
		common.ErrorWriter404(w, vars.URLParam("node"))
	}

	b, err := ioutil.ReadAll(r.Body)
	common.ErrorWriter(w, err)

	var x types.TaskInfo

	err = json.Unmarshal(b, &x)
	common.ErrorWriter(w, err)

	if x.ID == "hello" {
		if n, ok = node.GetNodeName(n); !ok {
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

func initQueues(tasks []types.Task) error {
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

		if err = queueTask(t.Node, t.ID, cmd); err != nil {
			return err
		}
	}

	return nil
}

func initAgent(agent string) error {
	refreshDB()

	var t []types.Task
	if err := db.QueryRow(context.Background(),
		"SELECT jsonb_agg(to_jsonb(tasks) - 'id_short') FROM tasks WHERE state = $1 AND node = $2;",
		types.StateQueued, agent).Scan(&t); err != nil {
		return err
	}

	if err := initQueues(t); err != nil {
		return err
	}

	return nil
}

// Init initializes the agent queues.
func Init() error {
	refreshDB()

	var t []types.Task
	if err := db.QueryRow(context.Background(),
		"SELECT jsonb_agg(to_jsonb(tasks) - 'id_short') FROM tasks WHERE state = $1;",
		types.StateQueued).Scan(&t); err != nil {
		return err
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
