package job

import (
	"reflect"
	"strings"
	"testing"

	"github.com/chabad360/covey/test"
)

func TestJobNew(t *testing.T) {
	//revive:disable:line-length-limit
	var tests = []struct {
		body string
		want string
	}{
		{`{"name":"test","nodes": ["test"],"tasks": {"update": {"plugin": "test","details": {"command": "test"}}}}`,
			`{"name":"test","id":"240875a9cf2c26d484a78b3f7f5aad21dd8f6e74031a7a5669f787d33e1b4cda","nodes":["test"],"tasks":{"update":{"plugin":"test","details":{"command":"test"}}},"task_history":[]}
`},
		{`{"name":"test","nodes": ["test"],"tasks": {"update": {"plugin": "test","details": {"command": "test"}}}}`,
			`{"error":"duplicate job: test"}
`},
		{`{"name":}`,
			`{"error":"models.Job.Name: ReadString: expects \" or n, but found }, error found in #9 byte of ...|{\"name\":}|..., bigger context ...|{\"name\":}|..."}
`},
	}
	//revive:enable:line-length-limit

	h := test.PureBoilerplate("GET", "/api/v1/jobs/new", jobNew)

	for _, tt := range tests {
		testname := tt.body
		t.Run(testname, func(t *testing.T) {
			rr, req, err := test.HTTPBoilerplate("GET", "/api/v1/jobs/new", strings.NewReader(tt.body))
			if err != nil {
				t.Fatal(err)
			}

			h.ServeHTTP(rr, req)
			if !reflect.DeepEqual(rr.Body.Bytes()[0:15], []byte(`{"name":"test",`)) && rr.Body.String() != tt.want {
				t.Errorf("jobNew body = %v, want %v", rr.Body.String(), tt.want)
			}
		})
	}
}

func TestJobGet(t *testing.T) {
	//revive:disable:line-length-limit
	var tests = []struct {
		id   string
		want string
	}{
		{"3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e",
			`{"name":"update","id":"3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e","nodes":["node1"],"tasks":{"update":{"plugin":"shell","details":{"command":"sudo apt update \u0026\u0026 sudo apt upgrade -y"}}}}
`},
		{"3", `{"error":"404 3 not found"}
`},
	}
	//revive:enable:line-length-limit

	h := test.PureBoilerplate("GET", "/api/v1/job/:job", jobGet)

	for _, tt := range tests {
		testname := tt.id
		t.Run(testname, func(t *testing.T) {
			rr, req, err := test.HTTPBoilerplate("GET", "/api/v1/job/"+tt.id, nil)
			if err != nil {
				t.Fatal(err)
			}

			h.ServeHTTP(rr, req)
			if rr.Body.String() != tt.want {
				t.Errorf("jobGet body = %v, want %v", rr.Body.String(), tt.want)
			}
		})
	}
}
