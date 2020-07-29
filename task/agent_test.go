package task

import (
	"github.com/chabad360/covey/test"
	"reflect"
	"strings"
	"testing"
)

func Test_queueTask(t *testing.T) {
	tests := []struct {
		name    string
		nodeID  string
		wantErr bool
		want    agentTask
	}{
		{"success", "3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e", false, agentTask{"test", "test"}},
		{"second", "3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e", false, agentTask{"test2", "test2"}},
		{"fail", "sadf", true, agentTask{"test", "test"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := queueTask(tt.nodeID, tt.want.ID, tt.want.Command); (err != nil) != tt.wantErr {
				t.Errorf("queueTask() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if task := queues[tt.nodeID].Back().Value.(agentTask); !reflect.DeepEqual(task, tt.want) {
					t.Errorf("queueTask() failed to add task, got %v, want %v", task, tt.want)
				}
			}
		})
	}
}

func Test_agentPost(t *testing.T) {
	queues = make(map[string]*List)
	tests := []struct {
		name string
		id   string
		body string
		want string
	}{
		{"success", "3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e", `null`, `{"0":{"command":"test1","id":"test1"}}
`},
		{"empty", "3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e", `null`, `null
`},
		{"send", "3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e", `{"id":"` + task1.ID + `", "log":["test"], "exit_code":0}`, `null
`},
		{"sendFail", "3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e", `{"id":"test2", "log":["test"], "exit_code":0}`, `{"error":"saveTask: task test2 not found"}
`},
		{"fail", "3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e", ``, `{"error":"readObjectStart: expect { or n, but found \u0000, error found in #0 byte of ...||..., bigger context ...||..."}
`},
		{"fail404", "3", ``, `{"error":"404 3 not found"}
`},
	}

	h := test.PureBoilerplate("POST", "/agent/:node", agentPost)
	queueTask("3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e", "test1", "test1")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr, req, err := test.HTTPBoilerplate("POST", "/agent/"+tt.id, strings.NewReader(tt.body))
			if err != nil {
				t.Fatal(err)
			}

			h.ServeHTTP(rr, req)
			if !reflect.DeepEqual(rr.Body.Bytes(), []byte(tt.want)) && rr.Body.String() != tt.want {
				t.Errorf("agentPost body = %v, want %v", rr.Body.String(), tt.want)
			}
		})
	}
}
