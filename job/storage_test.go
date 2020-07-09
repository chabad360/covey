package job

import (
	"context"
	"github.com/chabad360/covey/models"
	"log"
	"os"
	"testing"

	"github.com/chabad360/covey/storage"
	"github.com/chabad360/covey/test"
)

var j = &models.Job{
	Name:  "update",
	ID:    "3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e",
	Nodes: []string{"node1"},
	Tasks: map[string]models.JobTask{
		"update": {
			Plugin:  "shell",
			Details: map[string]string{"command": "sudo apt update && sudo apt upgrade -y"},
		},
	},
	TaskHistory: []string{},
}

func TestAddJob(t *testing.T) {
	//revive:disable:line-length-limit
	var tests = []struct {
		id   string
		want string
	}{
		{"3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e",
			`{"id": "3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e", "cron": "", "name": "update", "nodes": ["node1"], "tasks": {"update": {"plugin": "shell", "details": {"command": "sudo apt update && sudo apt upgrade -y"}}}, "task_history": []}`},
		{"3", ""},
	}
	//revive:enable:line-length-limit

	testError := AddJob(*j)

	for _, tt := range tests {
		testname := tt.id
		t.Run(testname, func(t *testing.T) {
			var got []byte
			if db.QueryRow(context.Background(), "SELECT to_jsonb(jobs) - 'id_short' FROM jobs WHERE id = $1;",
				tt.id).Scan(&got); string(got) != tt.want {
				t.Errorf("AddJob() = %v, want %v, error: %v", string(got), tt.want, testError)
			}
		})
	}
}

func TestUpdateJob(t *testing.T) {
	//revive:disable:line-length-limit
	var tests = []struct {
		id   string
		want string
	}{
		{"3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e",
			`{"id": "3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e", "cron": "5 * * * *", "name": "update", "nodes": ["node1"], "tasks": {"update": {"plugin": "shell", "details": {"command": "sudo apt update && sudo apt upgrade -y"}}}, "task_history": []}`},
		{"3", ""},
	}
	//revive:enable:line-length-limit

	ju := *j
	ju.Cron = "5 * * * *"
	testError := UpdateJob(ju)

	for _, tt := range tests {
		testname := tt.id
		t.Run(testname, func(t *testing.T) {
			var got []byte
			if db.QueryRow(context.Background(), "SELECT to_jsonb(jobs) - 'id_short' FROM jobs WHERE id = $1;",
				tt.id).Scan(&got); string(got) != tt.want {
				t.Errorf("UpdateJob() = %v, want %v, error: %v", string(got), tt.want, testError)
			}
		})
	}
}

func TestGetJobWithFullHistory(t *testing.T) {
	//revive:disable:line-length-limit
	var tests = []struct {
		id   string
		want string
	}{
		{"3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e",
			`{"id": "3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e", "cron": "5 * * * *", "name": "update", "nodes": ["node1"], "tasks": {"update": {"plugin": "shell", "details": {"command": "sudo apt update && sudo apt upgrade -y"}}}, "task_history": null}`},
		{"3", ""},
	}
	//revive:enable:line-length-limit

	for _, tt := range tests {
		testname := tt.id
		t.Run(testname, func(t *testing.T) {
			if got, err := GetJobWithFullHistory(tt.id); string(got) != tt.want {
				t.Errorf("GetJobWithFullHistory() = %v, want %v, error: %v", string(got), tt.want, err)
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

	db.Exec(context.Background(), `INSERT INTO jobs(id, id_short, name, cron, nodes, tasks, task_history)
	VALUES($1, $2, $3, $4, $5, $6, $7);`,
		j.ID, j.GetIDShort(), j.Name, j.Cron, j.Nodes, j.Tasks, j.TaskHistory)

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
