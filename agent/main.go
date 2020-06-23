package main

import (
	"bytes"
	"container/list"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/joho/godotenv"
	json "github.com/json-iterator/go"
)

const (
	sleepDuration = time.Second
)

var (
	queue       = list.New()
	buffer      = new(bytes.Buffer)
	currentTask task
	activeTask  *runningTask
	agent       *config
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
	log.Println("Agent ID:", agent.AgentID)
	log.Println("Host:", agent.Host)

	go runner()
	log.Println("Covey Agent started!")

	for {
		go everySecond()
		time.Sleep(sleepDuration)
	}
}

func settings(file string) (*config, error) {
	godotenv.Load(file)

	var conf = config{}
	var exists bool
	if conf.AgentID, exists = os.LookupEnv("AGENT_ID"); !exists || conf.AgentID == "" {
		return nil, fmt.Errorf("missing AGENT_ID")
	}

	if conf.LogLevel, exists = os.LookupEnv("LOG_LEVEL"); !exists || conf.LogLevel == "" {
		conf.LogLevel = "INFO"
	}

	var host string
	var port string
	if host, exists = os.LookupEnv("AGENT_HOST"); !exists || host == "" {
		return nil, fmt.Errorf("missing AGENT_HOST")
	}
	if port, exists = os.LookupEnv("AGENT_HOST_POST"); !exists || port == "" {
		port = "8080"
	}
	conf.Host = host + ":" + port
	agentPath = conf.Host + "/agent/" + conf.AgentID

	return &conf, nil
}

func everySecond() {
	// First read it, then delete it.
	body, err := json.Marshal(activeTask)
	errC(err)
	activeTask.Log = nil // Once we've read once, we don't want to read it again.
	if string(body) == "null" {
		body = nil
	}
	if activeTask.ExitCode != 257 {
		activeTask = nil
	}
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

// TODO: lower complexity levels of this code.
func runner() {
	for {
		if qt := queue.Front(); qt != nil {
			e := make(chan int)
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
					for b, err := buffer.ReadBytes('\n'); ; b, err = buffer.ReadBytes('\n') {
						if string(b) == "" {
							break
						}
						activeTask.Log = append(activeTask.Log, string(b))
						if err != nil {
							break
						}
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
