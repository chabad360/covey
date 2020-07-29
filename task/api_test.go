package task

import (
	"github.com/chabad360/covey/storage"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/chabad360/covey/test"
)

func TestTaskNew(t *testing.T) {
	var tests = []struct {
		name string
		body string
		want string
	}{
		// revive:disable:line-length-limit
		//		{"regular", `{"name":"update","nodes": ["test"],"tasks": {"update": {"plugin": "test","details": {"command": "test"}}}}`,
		//			`{"name":"update","id":"240875a9cf2c26d484a78b3f7f5aad21dd8f6e74031a7a5669f787d33e1b4cda","nodes":["test"],"tasks":{"update":{"plugin":"test","details":{"command":"test"}}},"task_history":[]}
		//`},
		//		{"cron", `{"name":"cron", "cron": "5 * * * *","nodes": ["test"],"tasks": {"update": {"plugin": "test","details": {"command": "test"}}}}`,
		//			`{"name":"cron","id":"240875a9cf2c26d484a78b3f7f5aad21dd8f6e74031a7a5669f787d33e1b4cda","cron":"5 * * * *","nodes":["test"],"tasks":{"update":{"plugin":"test","details":{"command":"test"}}},"task_history":[]}
		//`},
		{"error", `{"plugin":}`,
			`{"error":"models.Task.Plugin: ReadString: expects \" or n, but found }, error found in #9 byte of ...|{\"plugin\":}|..., bigger context ...|{\"plugin\":}|..."}
`},
		// revive:enable:line-length-limit
	}

	h := test.PureBoilerplate("POST", "/api/v1/tasks", taskNew)

	for _, tt := range tests {
		testname := tt.name
		t.Run(testname, func(t *testing.T) {
			rr, req, err := test.HTTPBoilerplate("POST", "/api/v1/tasks", strings.NewReader(tt.body))
			if err != nil {
				t.Fatal(err)
			}

			h.ServeHTTP(rr, req)
			if !reflect.DeepEqual(rr.Body.Bytes()[0:10], []byte(tt.want)[0:10]) && rr.Body.String() != tt.want {
				t.Errorf("taskNew body = %v, want %v", rr.Body.String(), tt.want)
			}
		})
	}
}

func TestTaskGet(t *testing.T) {
	var tests = []struct {
		name string
		id   string
		want string
	}{
		// revive:disable:line-length-limit
		{"success", task1.ID,
			`{"state":6,"plugin":"test","id":"91daa4d64a2693c0e9d012650b19e16c9f64541f8c34e24e4c387a4a8a44cb38","node":"test","details":{"test":"test"},"exit_code":258,"created_at":"2020-07-29T14:15:21.399917-07:00","updated_at":"2020-07-29T14:15:21.399917-07:00"}
`},
		{"fail", "3", `{"error":"404 3 not found"}
`},
		// revive:enable:line-length-limit
	}

	h := test.PureBoilerplate("GET", "/api/v1/task/:task", taskGet)

	for _, tt := range tests {
		testname := tt.name
		t.Run(testname, func(t *testing.T) {
			rr, req, err := test.HTTPBoilerplate("GET", "/api/v1/task/"+tt.id, nil)
			if err != nil {
				t.Fatal(err)
			}

			h.ServeHTTP(rr, req)
			if !reflect.DeepEqual(rr.Body.Bytes()[0:10], []byte(tt.want)[0:10]) && rr.Body.String() != tt.want {
				t.Errorf("taskGet body = %v, want %v", rr.Body.String(), tt.want)
			}
		})
	}
}

func TestMain(m *testing.M) {
	pool, resource, pdb, err := test.Boilerplate()
	storage.DB = pdb
	if err != nil {
		log.Fatalf("Could not setup DB connection: %s", err)
	}

	storage.AddTask(task1)

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
