package task

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/chabad360/covey/storage"
	"github.com/chabad360/covey/task/types"
	"github.com/chabad360/covey/test"
)

var task = &types.Task{
	ID:     "3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e",
	State:  types.StateRunning,
	Plugin: "test",
	Node:   "test",
	Time:   time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC),
	Details: struct {
		Test  string `json:"Test"`
		Test2 string `json:"Test2"`
	}{"test", "test"},
}

func TestAddTask(t *testing.T) {
	//revive:disable:line-length-limit
	var tests = []struct {
		id   string
		want string
	}{
		{"3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e",
			`{"id": "3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e", "log": null, "node": "test", "time": "2000-01-01T01:01:01.000000001Z", "state": 2, "plugin": "test", "details": {"Test": "test", "Test2": "test"}}`},
		{"3", ""},
	}
	//revive:enable:line-length-limit

	testError := addTask(task)

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.id)
		t.Run(testname, func(t *testing.T) {
			var got []byte
			if db.QueryRow(context.Background(), "SELECT to_jsonb(tasks) - 'id_short' FROM tasks WHERE id = $1;",
				tt.id).Scan(&got); string(got) != tt.want {
				t.Errorf("addTask() = %v, want %v, error: %v", string(got), tt.want, testError)
			}
		})
	}
}

func TestUpdateTask(t *testing.T) {
	//revive:disable:line-length-limit
	var tests = []struct {
		id   string
		want string
	}{
		{"3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e",
			`{"id": "3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e", "log": ["hello", "world"], "node": "test", "time": "2000-01-01T01:01:01.000000001Z", "state": 2, "plugin": "test", "details": {"Test": "test", "Test2": "test"}}`},
		{"3", ""},
	}
	//revive:enable:line-length-limit

	tu := task
	tu.Log = []string{"hello", "world"}
	testError := updateTask(tu)

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.id)
		t.Run(testname, func(t *testing.T) {
			var got []byte
			if db.QueryRow(context.Background(), "SELECT to_jsonb(tasks) - 'id_short' FROM tasks WHERE id = $1;",
				tt.id).Scan(&got); string(got) != tt.want {
				t.Errorf("updateTask() = %v, want %v, error: %v", string(got), tt.want, testError)
			}
		})
	}
}

func TestGetTaskJSON(t *testing.T) {
	//revive:disable:line-length-limit
	var tests = []struct {
		id   string
		want string
	}{
		{"3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e",
			`{"id": "3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e", "log": ["hello", "world"], "node": "test", "time": "2000-01-01T01:01:01.000000001Z", "state": 2, "plugin": "test", "details": {"Test": "test", "Test2": "test"}}`},
		{"3", ""},
	}
	//revive:enable:line-length-limit

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.id)
		t.Run(testname, func(t *testing.T) {
			if got, err := getTaskJSON(tt.id); string(got) != tt.want {
				t.Errorf("getTaskJSON() = %v, want %v, error: %v", string(got), tt.want, err)
			}
		})
	}
}

func TestMain(m *testing.M) {
	pool, resource, pdb, err := test.Boilerplate()
	db = pdb
	storage.DB = pdb
	if err != nil {
		log.Fatalf("Could not setup DB connection: %s", err)
	}

	db.Exec(context.Background(), `INSERT INTO tasks(id, id_short, plugin, state, node, time, log, details) 
		VALUES($1, $2, $3, $4, $5, $6, $7, $8);`,
		task.GetID(), task.GetIDShort(), task.GetPlugin(), task.GetState(), task.GetNode(),
		func() string { t, _ := task.GetTime().MarshalText(); return string(t) }(),
		task.GetLog(), task.GetDetails())

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
