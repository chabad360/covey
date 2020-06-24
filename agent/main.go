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
	queue       = tList{}
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

type tList struct{ list.List }

func (l *tList) UnmarshalJSON(b []byte) error {
	var m map[int]task
	err := json.Unmarshal(b, &m)
	if err != nil {
		return err
	}

	for i := 0; i < len(m); i++ {
		l.PushBack(m[i])
	}
	return nil
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
	agentPath = "http://" + conf.Host + "/agent/" + conf.AgentID

	return &conf, nil
}

func everySecond() {
	// First read it, then delete it.
	var body = []byte("{}")
	var err error
	if activeTask != nil {
		body, err = json.Marshal(activeTask)
		errC(err)
		activeTask.Log = nil // Once we've read once, we don't want to read it again.
		if activeTask.ExitCode != 257 {
			activeTask = nil
		}
	}
	r, err := http.Post(agentPath, "application/json", strings.NewReader(string(body)))
	errC(err)

	taskJSON, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(taskJSON, &queue)
	errC(err)

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
			var buffer bytes.Buffer
			e := make(chan int)
			t := qt.Value.(task)

			log.Printf(t.Command)
			cmd := exec.Command("/bin/bash", "-c", t.Command)
			cmd.Stdout = &buffer
			cmd.Stderr = &buffer

			activeTask = &runningTask{
				ID:       t.ID,
				ExitCode: 257,
			}
			err := cmd.Start()
			errC(err)

			go func() {
				err = cmd.Wait()
				if err != nil {
					if err, ok := err.(*exec.ExitError); ok {
						e <- err.ExitCode()
					}
				} else {
					e <- 0
				}
				close(e)
			}()

			for {
				select {
				case i := <-e:
					activeTask.ExitCode = i
					goto end
				default:
					// There is an issue where this method fails to capture anything after \r (until the next \n),
					// please help. Also, random nil pointer dereferences...
					for b, _ := buffer.ReadBytes('\n'); string(b) != ""; b, _ = buffer.ReadBytes('\n') {
						if b[len(b)-1] != 0 {
							if b[len(b)-1] == '\n' {
								activeTask.Log = append(activeTask.Log, string(b[:len(b)-1]))
							} else {
								activeTask.Log = append(activeTask.Log, string(b[:]))
							}
						}
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
