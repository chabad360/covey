package main

//go:generate go build -ldflags "-s -w" -trimpath -o ../assets/agent/agent github.com/chabad360/covey/agent
//go:generate upx ../assets/agent/agent

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"time"

	json "github.com/json-iterator/go"
)

func main() {
	var (
		err        error
		activeTask *runningTask
		agent      config
		q          = &queue{}
	)

	t := time.NewTicker(time.Second)
	rt := make(chan *runningTask, 1)

	if err = settings(&agent); err != nil {
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

	go taskManager(q, rt)

	for {
		<-t.C
		activeTask, err = everySecond(q, rt, activeTask, agent)
		if err != nil {
			log.Println(err)
		}
	}
}

func everySecond(q *queue, rt <-chan *runningTask, at *runningTask, agent config) (*runningTask, error) {
	if at == nil && len(rt) != 0 {
		at = <-rt
	}

	at, body, err := genBody(at)
	if err != nil {
		return at, err
	}

	got, err := connect(body, agent.AgentPath)
	if err != nil {
		return at, err
	}

	if err = json.Unmarshal(got, &q); err != nil {
		return at, err
	}

	return at, nil
}

func genBody(rt *runningTask) (*runningTask, []byte, error) {
	var t *returnTask

	if rt == nil {
		goto done
	}

	t = &returnTask{
		ID:  rt.ID,
		Log: rt.GetLog(),
	}

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

func connect(body []byte, path string) ([]byte, error) {
	var err error
	var r *http.Response
	for {
		r, err = http.Post(path, "application/json", bytes.NewReader(body)) //nolint:gosec
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

	cmd := exec.Command("/bin/sh", "-c", t.Command) //nolint:gosec

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		s = 11
		t.Log(fmt.Sprintf("agent error: %v", err))
		goto done
	}

	if err = cmd.Start(); err != nil {
		s = 11
		t.Log(fmt.Sprintf("agent error: %v", err))
		goto done
	}

	bb = bufio.NewScanner(stdout)
	for bb.Scan() {
		t.Log(bb.Text())
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
