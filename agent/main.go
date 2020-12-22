package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"sync"
	"time"

	json "github.com/json-iterator/go"
)

var (
	activeTask runningTask
	agent      config
	q          = &Queue{mutex: &sync.Mutex{}}
	timer      = time.NewTicker(time.Second)
	nextTask   baseContext
)

func main() {
	if err := settings(&agent); err != nil {
		log.Fatal(err)
	}

	// ignoring log level for now
	log.Println("Agent ID:", agent.ID)
	log.Println("Path:", agent.AgentPath)

	log.Println("Covey Agent started!")

	helloTask := newRunningTask(task{
		ID:      "hello",
		Command: "hello",
	})
	helloTask.Log("hello")
	helloTask.Finish(0, 0)
	activeTask = helloTask

	go taskManager(q)

	for {
		<-timer.C
		if err := everySecond(q, agent); err != nil {
			log.Println(err)
		}
	}
}

func everySecond(q *Queue, agent config) error {
	defer fmt.Println("sec")
	body, err := genBody(activeTask)
	if err != nil {
		return err
	}

	got, err := connect(body, agent.AgentPath)
	if err != nil {
		return err
	}

	return json.Unmarshal(got, &q)
}

func genBody(rt runningTask) ([]byte, error) {
	t := returnTask{
		ID:  rt.ID,
		Log: rt.GetLog(),
	}

	select {
	case <-nextTask.Done():
		t = returnTask{}
		return json.Marshal(t)
	case <-rt.context.Done():
		t.ExitCode = <-rt.ExitCode
		t.State = <-rt.State
		nextTask.Close()
	default:
		t.ExitCode = 257
		t.State = 2
	}

	return json.Marshal(t)
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

func taskManager(q *Queue) {
	for {
		qt := q.GetNext()        // GetNext the next task in the Queue
		t := newRunningTask(qt)  // Create a runningTask
		go run(t)                // Start the task
		nextTask = baseContext{} // Prevent the next task from running until this one is processed
		activeTask = t           // Set this task as the activeTask
		<-nextTask.Done()        // Wait for this task to finish processing before moving on
		//activeTask = runningTask{}
	}
}

func run(t runningTask) {
	cmd := exec.Command("/bin/sh", "-c", t.Command) //nolint:gosec

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Log(fmt.Sprintf("agent error: %v", err))
		t.Finish(0, 11)
		return
	}

	err = cmd.Start()
	if err != nil {
		t.Log(fmt.Sprintf("agent error: %v", err))
		t.Finish(0, 11)
		return
	}

	bb := bufio.NewScanner(stdout)
	for bb.Scan() {
		t.Log(bb.Text())
	}

	err = cmd.Wait()
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			t.Finish(e.ExitCode(), 1)
			return
		}
	}
	t.Finish(0, 0)
}
