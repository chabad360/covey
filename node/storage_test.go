package node

import (
	"github.com/chabad360/covey/models"
	"log"
	"os"
	"reflect"
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
	IP:         "127.0.0.1",
}

func TestAddNode(t *testing.T) {
	//revive:disable:line-length-limit
	var tests = []struct {
		id   string
		want *models.Node
	}{
		{"3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e", n},
		{"3", &models.Node{}},
	}
	//revive:enable:line-length-limit

	testError := addNode(n)

	for _, tt := range tests {
		testname := tt.id
		t.Run(testname, func(t *testing.T) {
			var got models.Node
			if db.Where("id = ?", tt.id).First(&got); reflect.DeepEqual(got, tt.want) {
				t.Errorf("addNode() = %v, want %v, error: %v", got, tt.want, testError)
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

	err = db.AutoMigrate(&models.Node{})
	if err != nil {
		log.Fatalf("error preping the database: %s", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
