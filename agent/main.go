package main

import (
	"bytes"
	"container/list"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const (
	sleepDuration = time.Second
)

var (
	queue       = list.New()
	buffer      = new(bytes.Buffer)
	currentTask task
	activeTask  *runningTask
	agent       config
	agentPath   string
)

// Task info
type runningTask struct {
	Log      []string `json:"log"`
	ExitCode int      `json:"exit_code"`
	ID       string   `json:"id"`
}

type task struct {
	Command string `json:"command"`
	ID      string `json:"id"`
}

type config struct {
	AgentID  string
	LogLevel string
	Host     string
}

func main() {
	var err error
	agent, err = settings("/etc/covey/agent.conf")
	errC(err)
	// ignoring log level for now

	for {
		go everySecond()
		time.Sleep(sleepDuration)
	}
}

func settings(file string) (config, error) {
	err := godotenv.Load(file)
	errC(err)

	var conf = config{}
	var exists bool
	if conf.AgentID, exists = os.LookupEnv("AGENT_ID"); !exists || conf.AgentID == "" {
		panic(fmt.Errorf("missing AGENT_ID"))
	}

	if conf.LogLevel, exists = os.LookupEnv("LOG_LEVEL"); !exists || conf.LogLevel == "" {
		conf.LogLevel = "INFO"
	}

	var host string
	var port string
	if host, exists = os.LookupEnv("AGENT_HOST"); !exists || host == "" {
		panic(fmt.Errorf("missing AGENT_HOST"))
	}
	if port, exists = os.LookupEnv("AGENT_HOST_POST"); !exists || port == "" {
		port = "8080"
	}
	conf.Host = host + ":" + port

	return conf, nil
}

func everySecond() {
	body, err := json.Marshal(activeTask)
	if string(body) == "null" {
		body = nil
	}
	if activeTask.ExitCode != 257 {
		activeTask = nil
	}
	activeTask.Log = nil
	errC(err)
	r, err := http.Post(agentPath, "application/json", strings.NewReader(string(body)))
	errC(err)

	taskJSON, err := ioutil.ReadAll(r.Body)
	if string(taskJSON) != "" {
		var newTask task
		err = json.Unmarshal(taskJSON, &newTask)
		errC(err)
		queue.PushBack(newTask)
	}
	r.Body.Close()
}

func errC(err error) {
	if err != nil {
		panic(err)
	}
}

func runner() {
	for {
		if queue.Front() != nil {
			e := make(chan int)
			qt := queue.Front()
			t := qt.Value.(task)

			cmd := exec.Command("/bin/sh", t.Command)
			cmd.Stdout = buffer
			cmd.Stderr = buffer
			activeTask = &runningTask{
				ID:       t.ID,
				ExitCode: 257,
			}
			cmd.Start()

			go func() {
				err := cmd.Wait()
				if err != nil {
					if err, ok := err.(*exec.ExitError); ok {
						e <- err.ExitCode()
					}
				} else {
					e <- 0
				}
				close(e)
			}()

			// This select loop is necessary because the buffer may occasionally be empty/missing \n (echo -n),
			// this is my way of handleing that race condition, plus echo -n (it doesn't pin a \n to the end of Stdout).
			for {
				select {
				case i := <-e:
					for b, _ := buffer.ReadBytes('\n'); ; b, _ = buffer.ReadBytes('\n') {
						if string(b) == "" {
							break
						}
						activeTask.Log = append(activeTask.Log, string(b))
					}
					activeTask.ExitCode = i
					goto end
				default:
					for b, err := buffer.ReadBytes('\n'); err == nil; b, err = buffer.ReadBytes('\n') {
						activeTask.Log = append(activeTask.Log, string(b))
					}
				}
			}

		end:
			buffer.Reset()
			queue.Remove(qt)
		}
		time.Sleep(sleepDuration * 2) // Two seconds between each task
	}
}