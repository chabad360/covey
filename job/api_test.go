package job

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/chabad360/covey/models"
	"github.com/chabad360/covey/storage"
	"github.com/chabad360/covey/test"
)

var (
	j  = test.J1
	j2 = test.J2
)

func TestJobNew(t *testing.T) {
	var tests = []struct {
		name string
		body string
		want string
	}{
		// revive:disable:line-length-limit
		{"regular", `{"name":"update","nodes": ["test"],"tasks": {"update": {"plugin": "test","details": {"command": "test"}}}}`,
			`{"name":"update","id":"240875a9cf2c26d484a78b3f7f5aad21dd8f6e74031a7a5669f787d33e1b4cda","nodes":["test"],"tasks":{"update":{"plugin":"test","details":{"command":"test"}}},"task_history":[]}
`},
		{"cron", `{"name":"cron", "cron": "5 * * * *","nodes": ["test"],"tasks": {"update": {"plugin": "test","details": {"command": "test"}}}}`,
			`{"name":"cron","id":"240875a9cf2c26d484a78b3f7f5aad21dd8f6e74031a7a5669f787d33e1b4cda","cron":"5 * * * *","nodes":["test"],"tasks":{"update":{"plugin":"test","details":{"command":"test"}}},"task_history":[]}
`},
		{"errorDuplicate", `{"name":"update","nodes": ["test"],"tasks": {"update": {"plugin": "test","details": {"command": "test"}}}}`,
			`{"error":"duplicate job: update"}
`},
		{"error", `{"name":}`,
			`{"error":"models.Job.Name: ReadString: expects \" or n, but found }, error found in #9 byte of ...|{\"name\":}|..., bigger context ...|{\"name\":}|..."}
`},
		// revive:enable:line-length-limit
	}

	h := test.PureBoilerplate("POST", "/api/v1/jobs", jobNew)

	for _, tt := range tests {
		testname := tt.name
		t.Run(testname, func(t *testing.T) {
			rr, req, err := test.HTTPBoilerplate("POST", "/api/v1/jobs", strings.NewReader(tt.body))
			if err != nil {
				t.Fatal(err)
			}

			h.ServeHTTP(rr, req)
			if !cmp.Equal(rr.Body.Bytes()[0:10], []byte(tt.want)[0:10]) && rr.Body.String() != tt.want {
				t.Errorf("jobNew body = %v, want %v", rr.Body.String()[0:10], tt.want[0:10])
			}
		})
	}
}

func TestJobsGet(t *testing.T) {
	storage.DB.Delete(&models.Job{}, "id != ''")
	storage.AddJob(&j)
	storage.AddJob(&j2)
	js, _ := json.Marshal(j2)

	var tests = []struct {
		name   string
		params string
		want   string
	}{
		// revive:disable:line-length-limit
		{"success", "sortby=name", `["` + j2.ID + `","` + j.ID + `"]
`},
		{"onlyOne", "sortby=name&limit=1", `["` + j2.ID + `"]
`},
		{"offsetOne", "sortby=name&limit=1&offset=1", `["` + j.ID + `"]
`},
		{"expandOne", "sortby=name&limit=1&expand=true", `[` + string(js) + `]
`},
		// revive:enable:line-length-limit
	}

	h := test.PureBoilerplate("GET", "/api/v1/jobs", jobsGet)

	for _, tt := range tests {
		testname := tt.name
		t.Run(testname, func(t *testing.T) {
			rr, req, err := test.HTTPBoilerplate("GET", "/api/v1/jobs?"+tt.params, nil)
			if err != nil {
				t.Fatal(err)
			}

			h.ServeHTTP(rr, req)
			if !cmp.Equal(rr.Body.Bytes(), []byte(tt.want)) && rr.Body.String() != tt.want {
				t.Errorf("jobGet body = %v, want %v", rr.Body.String(), tt.want)
			}
		})
	}
}

func TestJobGet(t *testing.T) {
	storage.DB.Delete(&models.Job{}, "id != ''")
	storage.AddJob(&j2)
	var tests = []struct {
		name string
		id   string
		want string
	}{
		// revive:disable:line-length-limit
		{"success", "add",
			`{"name":"add","id":"3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e","nodes":["node1"],"tasks":{"update":{"plugin":"shell","details":{"command":"sudo apt update \u0026\u0026 sudo apt upgrade -y"}}}}
`},
		{"fail", "3", `{"error":"404 3 not found"}
`},
		// revive:enable:line-length-limit
	}

	h := test.PureBoilerplate("GET", "/api/v1/job/:job", jobGet)

	for _, tt := range tests {
		testname := tt.name
		t.Run(testname, func(t *testing.T) {
			rr, req, err := test.HTTPBoilerplate("GET", "/api/v1/job/"+tt.id, nil)
			if err != nil {
				t.Fatal(err)
			}

			h.ServeHTTP(rr, req)
			if !cmp.Equal(rr.Body.Bytes()[0:10], []byte(tt.want)[0:10]) && rr.Body.String() != tt.want {
				t.Errorf("jobGet body = %v, want %v", rr.Body.String()[0:10], tt.want[0:10])
			}
		})
	}
}

func TestJobUpdate(t *testing.T) {
	storage.DB.Delete(&models.Job{}, "id != ''")
	storage.AddJob(&j2)
	var tests = []struct {
		name string
		id   string
		body string
		want string
	}{
		// revive:disable:line-length-limit
		{"success", "add", `{"name":"add","cron":"5 * * * *","nodes": ["test"],"tasks": {"update": {"plugin": "test","details": {"command": "hello"}}}}`,
			`{"name":"add","id":"240875a9cf2c26d484a78b3f7f5aad21dd8f6e74031a7a5669f787d33e1b4cda","cron":"5 * * * *","nodes":["test"],"tasks":{"update":{"plugin":"test","details":{"command":"hello"}}},"task_history":[]}
`},
		{"error", "cron", `{"name":}`,
			`{"error":"models.Job.Name: ReadString: expects \" or n, but found }, error found in #9 byte of ...|{\"name\":}|..., bigger context ...|{\"name\":}|..."}
`},
		{"404", "c", "", `{"error":"404 c not found"}
`},
		// revive:enable:line-length-limit
	}

	h := test.PureBoilerplate("PUT", "/api/v1/jobs/:job", jobUpdate)

	for _, tt := range tests {
		testname := tt.name
		t.Run(testname, func(t *testing.T) {
			rr, req, err := test.HTTPBoilerplate("PUT", "/api/v1/jobs/"+tt.id, strings.NewReader(tt.body))
			if err != nil {
				t.Fatal(err)
			}

			h.ServeHTTP(rr, req)
			if !cmp.Equal(rr.Body.Bytes()[0:10], []byte(tt.want)[0:10]) && rr.Body.String() != tt.want {
				t.Errorf("jobNew body = %v, want %v", rr.Body.String(), tt.want)
			}
		})
	}
}

func TestJobDelete(t *testing.T) {
	storage.DB.Delete(&models.Job{}, "id != ''")
	storage.AddJob(&j2)
	var tests = []struct {
		name string
		id   string
		want string
	}{
		{"success", "add", `"add"
`},
		{"fail", "3", `{"error":"404 3 not found"}
`},
	}

	h := test.PureBoilerplate("DELETE", "/api/v1/job/:job", jobDelete)

	for _, tt := range tests {
		testname := tt.name
		t.Run(testname, func(t *testing.T) {
			rr, req, err := test.HTTPBoilerplate("DELETE", "/api/v1/job/"+tt.id, nil)
			if err != nil {
				t.Fatal(err)
			}

			h.ServeHTTP(rr, req)
			if rr.Body.String() != tt.want {
				t.Errorf("jobDelete body = %v, want %v", rr.Body.String(), tt.want)
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

	storage.AddJob(&j2)

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
