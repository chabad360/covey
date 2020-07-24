package storage

import (
	"github.com/chabad360/covey/models"
	"reflect"
	"testing"
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

	testError := AddNode(n)

	for _, tt := range tests {
		testname := tt.id
		t.Run(testname, func(t *testing.T) {
			var got models.Node
			if DB.Where("id = ?", tt.id).First(&got); reflect.DeepEqual(got, tt.want) {
				t.Errorf("addNode() = %v, want %v, error: %v", got, tt.want, testError)
			}
		})
	}
}

func TestGetNodeIDOrName(t *testing.T) {
	var tests = []struct {
		name  string
		id    string
		field string
		want  string
		want2 bool
	}{
		{"ok_ID", n.Name, "id", n.ID, true},
		{"notok_ID", "n", "id", "", false},
		{"ok_Name", n.ID, "name", n.Name, true},
		{"notok_Name", "n", "name", "", false},
	}

	for _, tt := range tests {
		testname := tt.name
		t.Run(testname, func(t *testing.T) {
			got, got2 := GetNodeIDorName(tt.id, tt.field)
			if got2 != tt.want2 {
				t.Errorf("GetNodeIDorName() = %v, want %v", got2, tt.want2)
			}
			if got != tt.want {
				t.Errorf("GetNodeIDorName() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func TestMain(m *testing.M) {
//	pool, resource, pdb, err := test.Boilerplate()
//	db = pdb
//	DB = pdb
//	if err != nil {
//		log.Fatalf("Could not setup DB connection: %s", err)
//	}
//
//	code := m.Run()
//
//	// You can't defer this because os.Exit doesn't care for defer
//	if err := pool.Purge(resource); err != nil {
//		log.Fatalf("Could not purge resource: %s", err)
//	}
//
//	os.Exit(code)
//}
