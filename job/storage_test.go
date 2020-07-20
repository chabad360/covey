package job

import (
	"github.com/chabad360/covey/models"
	"log"
	"os"
	"reflect"
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
}

func TestAddJob(t *testing.T) {
	//revive:disable:line-length-limit
	var tests = []struct {
		id   string
		want *models.Job
	}{
		{"3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e", j},
		{"3", &models.Job{}},
	}
	//revive:enable:line-length-limit

	testError := addJob(j)

	for _, tt := range tests {
		testname := tt.id
		t.Run(testname, func(t *testing.T) {
			var got models.Task
			if db.Where("id = ?", tt.id).First(&got); reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddJob() = %v, want %v, error: %v", got, tt.want, testError)
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

	j.Cron = "5 * * * *"
	testError := updateJob(j)

	for _, tt := range tests {
		testname := tt.id
		t.Run(testname, func(t *testing.T) {
			var got models.Job
			if db.Where("id = ?", tt.id).Scan(&got); reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateJob() = %v, want %v, error: %v", got, tt.want, testError)
			}
		})
	}
}

func TestGetJobWithFullHistory(t *testing.T) {
	//revive:disable:line-length-limit
	var tests = []struct {
		id   string
		want *models.Job
	}{
		{"3778ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e", j},
		{"3", &models.Job{}},
	}
	//revive:enable:line-length-limit

	for _, tt := range tests {
		testname := tt.id
		t.Run(testname, func(t *testing.T) {
			if got, err := getJobWithFullHistory(tt.id); reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetJobWithFullHistory() = %v, want %v, error: %v", got, tt.want, err)
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

	err = db.AutoMigrate(&models.Task{}, &models.Job{})
	if err != nil {
		log.Fatalf("error preping the database: %s", err)
	}

	db.Create(j)

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
