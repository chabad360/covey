package main

import (
	"bufio"
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
	activeTask      *runningTask
	agent           *config
	agentPath       string
	logChannel      = make(chan string, 1024)
	exitCodeChannel = make(chan int)
	taskIDChannel   = make(chan string)
	queue           = make(chan task, 1024)
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
		everySecond()
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
		var at *runningTask
		activeTask.Log = []string{<-logChannel}
		select {
		case e := <-exitCodeChannel:
			activeTask.ExitCode = e
			at = nil
		default:
			activeTask.ExitCode = 257
			at = activeTask
		}
		body, err = json.Marshal(activeTask)
		errC(err)
		activeTask = at
	} else {
		select {
		case t := <-taskIDChannel:
			activeTask = &runningTask{
				ID:       t,
				ExitCode: 257,
			}
		default:
			break
		}
	}
	r, err := http.Post(agentPath, "application/json", strings.NewReader(string(body)))
	errC(err)

	taskJSON, err := ioutil.ReadAll(r.Body)
	errC(err)
	var m map[int]task
	err = json.Unmarshal(taskJSON, &m)
	errC(err)

	for i := 0; i < len(m); i++ {
		queue <- m[i]
	}

	r.Body.Close()
}

func errC(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func runner() {
	for {
		t := <-queue

		taskIDChannel <- t.ID
		cmd := exec.Command("/bin/bash", "-c", t.Command)
		stdout, err := cmd.StdoutPipe()
		errC(err)

		err = cmd.Start()
		errC(err)

		bb := bufio.NewScanner(stdout)
		for bb.Scan() {
			logChannel <- bb.Text()
			log.Println(bb.Text())
		}
		err = cmd.Wait()
		if err != nil {
			if err, ok := err.(*exec.ExitError); ok {
				exitCodeChannel <- err.ExitCode()
			}
		} else {
			exitCodeChannel <- 0
		}
	}
}
