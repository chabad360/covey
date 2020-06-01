package task

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"plugin"

	"github.com/chabad360/covey/task/types"
)

func NewTask(taskJSON []byte) (types.ITask, error) {
	var task types.Task
	if err := json.Unmarshal(taskJSON, &task); err != nil {
		return nil, err
	}

	p, err := loadPlugin(task.Plugin)
	if err != nil {
		return nil, err
	}

	t, err := p.NewTask(taskJSON)
	if err != nil {
		return nil, err
	}

	saveConfig(t)
	return t, nil
}

// LoadConfig loads up the stored nodes
func LoadConfig() {
	log.Println("Loading Task Config")
	f, err := os.Open("./config/tasks.json")
	if err != nil {
		log.Println("Error loading task config")
		return
	}
	defer f.Close()

	var h map[string]json.RawMessage
	if err = json.NewDecoder(f).Decode(&h); err != nil {
		log.Fatal(err)
	}

	// Make this dynamic
	var plugins = make(map[string]types.TaskPlugin)
	p, err := loadPlugin("shell") // Hardcoding for now
	if err != nil {
		log.Fatal(err)
	}
	plugins["shell"] = p

	for _, node := range h {
		var z types.Task
		j, err := node.MarshalJSON()
		if err != nil {
			log.Fatal(err)
		}
		if err := json.Unmarshal(j, &z); err != nil {
			log.Fatal(err)
		}

		t, err := plugins[z.Plugin].LoadTask(j)
		if err != nil {
			log.Fatal(err)
		}

		tasks[t.GetID()] = t
		tasksShort[t.GetIDShort()] = t.GetID()
	}
}

func loadPlugin(pluginName string) (types.TaskPlugin, error) {
	p, err := plugin.Open("./plugins/task/" + pluginName + ".so")
	if err != nil {
		return nil, err
	}

	n, err := p.Lookup("Plugin")
	if err != nil {
		return nil, err
	}

	var s types.TaskPlugin
	s, ok := n.(types.TaskPlugin)
	if !ok {
		return nil, fmt.Errorf(pluginName, " does not provide a TaskPlugin")
	}

	return s, nil
}

func getTask(identifier string) (types.ITask, bool) {
	if t, ok := tasks[identifier]; ok {
		return t, true
	} else if t, ok := tasksShort[identifier]; ok {
		return tasks[t], true
	}
	return nil, false
}

func saveConfig(t types.ITask) {
	tasks[t.GetID()] = t
	tasksShort[t.GetIDShort()] = t.GetID()

	j, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create("./config/tasks.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err = f.Chmod(0600); err != nil {
		log.Fatal(err)
	}

	if _, err = f.Write(j); err != nil {
		log.Fatal(err)
	}
}
