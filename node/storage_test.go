package node

import (
	"context"
	"github.com/chabad360/covey/models"
	"log"
	"os"
	"testing"

	"github.com/chabad360/covey/storage"
	"github.com/chabad360/covey/test"
)

var n = &models.Node{
	Name:       "node",
	ID:         "3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e",
	PrivateKey: []byte("12345"),
	PublicKey:  []byte("12345"),
	HostKey:    []byte("12345"),
	Username:   "user",
}

func TestAddNode(t *testing.T) {
	//revive:disable:line-length-limit
	var tests = []struct {
		id   string
		want string
	}{
		{"3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e",
			`{"id": "3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e", "ip": "", "name": "node", "port": "", "host_key": "\\x3132333435", "username": "user", "public_key": "\\x3132333435", "private_key": "\\x3132333435"}`},
		{"3", ""},
	}
	//revive:enable:line-length-limit

	testError := addNode(n)

	for _, tt := range tests {
		testname := tt.id
		t.Run(testname, func(t *testing.T) {
			var got []byte
			if db.QueryRow(context.Background(), "SELECT to_jsonb(nodes) - 'id_short' FROM nodes WHERE id = $1;",
				tt.id).Scan(&got); string(got) != tt.want {
				t.Errorf("addNode() = %v, want %v, error: %v", string(got), tt.want, testError)
			}
		})
	}
}

func TestGetNodeID(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		if id, ok := GetNodeIDorName(n.Name, "id"); !ok && id != n.ID {
			t.Errorf("GetNodeID() = %v, want %v", id, n.ID)
		}
	})
	t.Run("not ok", func(t *testing.T) {
		if id, ok := GetNodeIDorName("n", "id"); ok && id == n.ID {
			t.Errorf("GetNodeID() = %v, want %v", id, n.ID)
		}
	})
}

func TestMain(m *testing.M) {
	pool, resource, pdb, err := test.Boilerplate()
	db = pdb
	storage.DB = pdb
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
