package main

import (
	"bytes"
	"container/list"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
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

// Config types
type config struct {
	Agent  agentConfig  `toml:"agent"`
	Client clientConfig `toml:"client"`
}

type agentConfig struct {
	AgentID  string `toml:"id"`
	LogLevel int    `toml:"log-level,omitempty"`
}

type clientConfig struct {
	Host     string `toml:"host"`
	HostPort string `toml:"host-port,omitempty"`
}

func main() {
	_, err := toml.DecodeFile("/etc/covey/agent.toml", &agent)
	errC(err)
	var hostPort string
	if hostPort = agent.Client.HostPort; hostPort == "" {
		hostPort = "8080"
	}
	host := agent.Client.Host + ":" + hostPort
	id := agent.Agent.AgentID
	agentPath = host + "/agent/" + id
	// ignoring log level for now

	for {
		go everySecond()
		time.Sleep(sleepDuration)
	}
}

func everySecond() {
	body, err := json.Marshal(activeTask)
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

			// This select loop is necessary because the buffer may occasionally be empty/missing \n,
			// so this is my way of handleing this race condition, and ensuring that a result from echo -n is captured.
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
		time.Sleep(sleepDuration * 2)
	}
}
