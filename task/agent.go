package task

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/node"
	"github.com/chabad360/covey/task/types"
	"github.com/go-playground/pure/v5"
)

var queues = make(map[string]*types.TaskList)

// QueueTask prepares a task to be sent to the node.
func QueueTask(nodeID string, taskID string, taskCommand string) error {
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
	vars := pure.RequestVars(r)
	n, ok := node.GetNodeID(vars.URLParam("node"))
	if !ok {
		common.ErrorWriter404(w, vars.URLParam("node"))
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		common.ErrorWriter(w, err)
	}
	log.Printf(string(b))
	var x types.TaskInfo
	err = json.Unmarshal(b, &x)
	if err != nil {
		common.ErrorWriter(w, err)
	}
	SaveTask(&x)

	common.Write(w, queues[n])
	delete(queues, n)
}

func RegisterAgentHandlers(r pure.IRouteGroup) {
	r.Post("/:node", agentPost)
}
