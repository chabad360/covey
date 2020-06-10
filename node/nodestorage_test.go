package node

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/chabad360/covey/node/types"
	"github.com/chabad360/covey/test"
)

var n = &types.Node{
	Name:    "node",
	ID:      "3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e",
	Details: struct{ Test string }{Test: "test"},
	Plugin:  "test",
}

func TestAddNode(t *testing.T) {
	//revive:disable:line-length-limit
	var tests = []struct {
		id   string
		want string
	}{
		{"3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e",
			`{"id": "3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e", "name": "node", "plugin": "test", "details": {"Test": "test"}}`},
		{"3", ""},
	}
	//revive:enable:line-length-limit

	testError := addNode(n)

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.id)
		t.Run(testname, func(t *testing.T) {
			var got []byte
			if db.QueryRow(context.Background(), "SELECT to_jsonb(nodes) - 'id_short' FROM nodes WHERE id = $1;",
				tt.id).Scan(&got); string(got) != tt.want {
				t.Errorf("addNode() = %v, want %v, error: %v", string(got), tt.want, testError)
			}
		})
	}
}

func TestMain(m *testing.M) {
	pool, resource, pdb, err := test.Boilerplate()
	db = pdb
	if err != nil {
		log.Fatalf("Could not setup DB connection: %s", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
