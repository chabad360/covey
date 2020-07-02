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
	activeID        string
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

	activeID = "hello"

	go func() {
		logChannel <- "hello"
		exitCodeChannel <- 1
	}()
	go runner()
	log.Println("Covey Agent started!")

	for {
		everySecond()
		time.Sleep(sleepDuration)
	}
}

func settings(file string) (*config, error) {
	var exists bool

	if err := godotenv.Load(file); err != nil {
		return nil, err
	}

	conf := config{}
	if conf.AgentID, exists = os.LookupEnv("AGENT_ID"); !exists || conf.AgentID == "" {
		return nil, fmt.Errorf("missing AGENT_ID")
	}

	if conf.LogLevel, exists = os.LookupEnv("LOG_LEVEL"); !exists || conf.LogLevel == "" {
		conf.LogLevel = "INFO"
	}

	var host string
	if host, exists = os.LookupEnv("AGENT_HOST"); !exists || host == "" {
		return nil, fmt.Errorf("missing AGENT_HOST")
	}

	var port string
	if port, exists = os.LookupEnv("AGENT_HOST_POST"); !exists || port == "" {
		port = "8080"
	}

	conf.Host = host + ":" + port
	agentPath = "http://" + conf.Host + "/agent/" + conf.AgentID

	return &conf, nil
}

func everySecond() {
	var err error

	body := []byte("{}")

	if activeID != "" {
		at := &runningTask{
			ID: activeID,
		}
		at.Log = []string{<-logChannel}

		select {
		case e := <-exitCodeChannel:
			at.ExitCode = e
			activeID = ""
		default:
			at.ExitCode = 257
		}

		body, err = json.Marshal(at)
		errC(err)
	} else {
		select {
		case t := <-taskIDChannel:
			activeID = t
		default:
			break
		}
	}

	var r *http.Response
	for {
		r, err = http.Post(agentPath, "application/json", strings.NewReader(string(body))) //nolint:gosec
		if err == nil {
			break
		}

		log.Printf("Couldn't connect to the host: %v\n", err)
		log.Println("Trying again in 5 seconds...")
		time.Sleep(4 * time.Second)
	}

	taskJSON, err := ioutil.ReadAll(r.Body)
	errC(err)

	var m map[int]task
	err = json.Unmarshal(taskJSON, &m)
	errC(err)

	for i := 0; i < len(m); i++ {
		queue <- m[i]
	}

	err = r.Body.Close()
	errC(err)
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
		cmd := exec.Command("/bin/bash", "-c", t.Command) //nolint:gosec
		stdout, err := cmd.StdoutPipe()
		errC(err)

		err = cmd.Start()
		errC(err)

		bb := bufio.NewScanner(stdout)
		for bb.Scan() {
			logChannel <- bb.Text()

			log.Println(bb.Text())
		}

		if err = cmd.Wait(); err != nil {
			if e, ok := err.(*exec.ExitError); ok {
				exitCodeChannel <- e.ExitCode()
			}
		} else {
			exitCodeChannel <- 0
		}
	}
}
