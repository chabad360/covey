package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"time"

	json "github.com/json-iterator/go"
)

var (
	agent      *config
	q          = queue{}
	activeTask *runningTask
)

func main() {
	var err error
	t := time.NewTicker(time.Second)
	rt := make(chan *runningTask, 1)
	q.Init()

	agent, err = settings()
	if err != nil {
		log.Fatal(err)
	}
	// ignoring log level for now
	log.Println("Agent ID:", agent.ID)
	log.Println("Path:", agent.AgentPath)

	log.Println("Covey Agent started!")

	ft := newRunningTask(task{
		ID:      "hello",
		Command: "hello",
	})
	ft.Finish(0, 0)
	rt <- ft

	go taskManager(&q, rt)

	for {
		<-t.C
		activeTask = everySecond(&q, rt, activeTask)
	}
}

func everySecond(q *queue, rt <-chan *runningTask, at *runningTask) *runningTask {
	if at == nil && len(rt) != 0 {
		at = <-rt
	}

	at, body, err := genBody(at)
	log.Println(err)

	got, err := connect(body)
	log.Println(err)

	err = json.Unmarshal(got, &q)
	log.Println(err)

	return at
}

func genBody(rt *runningTask) (*runningTask, []byte, error) {
	var t *returnTask

	if rt == nil {
		goto done
	}
	t = &returnTask{}

log:
	for {
		select {
		case l := <-rt.Log:
			t.Log = append(t.Log, l)
		default:
			break log
		}
	}

	t.ID = rt.ID

	select {
	case <-rt.Done():
		t.ExitCode = rt.ExitCode
		t.State = rt.State
		rt = nil
	default:
		t.ExitCode = 257
		t.State = 2
	}

done:
	b, err := json.Marshal(t)
	return rt, b, err
}

func connect(body []byte) ([]byte, error) {
	var err error
	var r *http.Response
	for {
		r, err = http.Post(agent.AgentPath, "application/json", bytes.NewReader(body)) //nolint:gosec
		if err == nil {
			break
		}

		log.Printf("Couldn't connect to the host: %v\n", err)
		log.Println("Trying again in 5 seconds...")
		time.Sleep(5 * time.Second)
	}
	defer r.Body.Close()

	return ioutil.ReadAll(r.Body)
}

func taskManager(q *queue, rt chan<- *runningTask) {
	for {
		qt := q.Get()
		log.Printf("new task %v", qt.ID)
		t := newRunningTask(qt)
		go run(t)
		rt <- t
		<-t.Done()
	}
}

func run(t *runningTask) {
	var bb *bufio.Scanner
	var ec int
	var s int

	cmd := exec.Command("/bin/bash", "-c", t.Cmd) //nolint:gosec
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		s = 11
		goto done
	}

	err = cmd.Start()
	if err != nil {
		s = 11
		goto done
	}

	bb = bufio.NewScanner(stdout)
	for bb.Scan() {
		t.Log <- bb.Text()
	}

	if err = cmd.Wait(); err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			ec = e.ExitCode()
			s = 1
		}
	}

done:
	t.Finish(ec, s)
}
