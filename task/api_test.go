package task

import (
	"log"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	json "github.com/json-iterator/go"

	"github.com/chabad360/covey/models"
	"github.com/chabad360/covey/storage"
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
			if !cmp.Equal(rr.Body.Bytes()[0:10], []byte(tt.want)[0:10]) && rr.Body.String() != tt.want {
				t.Errorf("taskNew body = %v, want %v", rr.Body.String(), tt.want)
			}
		})
	}
}

func TestTasksGet(t *testing.T) {
	storage.DB.Delete(&models.Task{}, "id != ''")
	storage.AddTask(&t1)
	storage.AddTask(&t2)
	js, _ := json.Marshal(t2)

	var tests = []struct {
		name   string
		params string
		want   string
	}{
		// revive:disable:line-length-limit
		{"success", "sortby=node", `["` + t2.ID + `","` + t1.ID + `"]
`},
		{"onlyOne", "sortby=node&limit=1", `["` + t2.ID + `"]
`},
		{"offsetOne", "sortby=node&limit=1&offset=1", `["` + t1.ID + `"]
`},
		{"expandOne", "sortby=node&limit=1&expand=true", `[` + string(js) + `]
`},
		// revive:enable:line-length-limit
	}

	h := test.PureBoilerplate("GET", "/api/v1/tasks", tasksGet)

	for _, tt := range tests {
		testname := tt.name
		t.Run(testname, func(t *testing.T) {
			rr, req, err := test.HTTPBoilerplate("GET", "/api/v1/tasks?"+tt.params, nil)
			if err != nil {
				t.Fatal(err)
			}

			h.ServeHTTP(rr, req)
			if !cmp.Equal(rr.Body.Bytes(), []byte(tt.want)) && rr.Body.String() != tt.want {
				t.Errorf("nodesGet body = %v, want %v", rr.Body.String(), tt.want)
			}
		})
	}
}

func TestTaskGet(t *testing.T) {
	js, _ := json.Marshal(t1)
	var tests = []struct {
		name string
		id   string
		want string
	}{
		// revive:disable:line-length-limit
		{"success", t1.ID, string(js) + "\n"},
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
			if rr.Body.String() != tt.want {
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

	storage.AddTask(&t1)
	storage.AddNode(&node1)

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
