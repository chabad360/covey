package job

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/chabad360/covey/common"
	"github.com/chabad360/covey/job/types"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
)

func TestGetJobWithFullHistory(t *testing.T) {
	var tests = []struct {
		id   string
		want string
	}{
		{"3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e",
			`{"id": "3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e", "cron": "", "name": "update", "nodes": ["node1"], "tasks": {"update": {"plugin": "shell", "details": {"Command": ["sudo apt update && sudo apt upgrade -y"]}}}, "task_history": null}`},
		{"3", ""},
	}

	j := &types.Job{
		Name:  "update",
		ID:    "3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e",
		Nodes: []string{"node1"},
		Tasks: map[string]types.JobTask{
			"update": {
				Plugin:  "shell",
				Details: struct{ Command []string }{Command: []string{"sudo apt update && sudo apt upgrade -y"}},
			},
		},
		TaskHistory: []string{},
	}

	db.Exec(`INSERT INTO jobs(id, id_short, name, cron, nodes, tasks, task_history)
		VALUES($1, $2, $3, $4, $5, $6, $7);`,
		j.ID, j.GetIDShort(), j.Name, j.Cron,
		common.UnsafeMarshal(j.Nodes), common.UnsafeMarshal(j.Tasks), common.UnsafeMarshal(j.TaskHistory))

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.id)
		t.Run(testname, func(t *testing.T) {
			if got, _ := GetJobWithFullHistory(tt.id); string(got) != tt.want {
				t.Errorf("GetJobWithFullHistory() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestMain(m *testing.M) {
	var err error
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("postgres", "9.6", []string{"POSTGRES_PASSWORD=secret", "POSTGRES_DB=covey"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err = pool.Retry(func() error {
		var err error
		db, err = sql.Open("postgres", fmt.Sprintf("postgres://postgres:secret@localhost:%s/%s?sslmode=disable",
			resource.GetPort("5432/tcp"), "covey"))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	db.Exec(`CREATE TABLE jobs (
		id TEXT PRIMARY KEY NOT NULL,
		id_short TEXT UNIQUE NOT NULL,
		name TEXT UNIQUE NOT NULL,
		cron TEXT,
		nodes JSONB NOT NULL,
		tasks JSONB NOT NULL,
		task_history JSONB
	);`)

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
