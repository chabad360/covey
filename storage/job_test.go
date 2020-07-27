package storage

import (
	"github.com/chabad360/covey/models"
	"reflect"
	"testing"
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

	testError := AddJob(j)

	for _, tt := range tests {
		testname := tt.id
		t.Run(testname, func(t *testing.T) {
			var got models.Task
			if DB.Where("id = ?", tt.id).First(&got); reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddJob() = %v, want %v, error: %v", got, tt.want, testError)
			}
		})
	}
}

func TestGetJob(t *testing.T) {
	//revive:disable:line-length-limit
	var tests = []struct {
		id    string
		want  *models.Job
		want2 bool
	}{
		{"update", j, true},
		{"3", &models.Job{}, false},
	}
	//revive:enable:line-length-limit

	for _, tt := range tests {
		testname := tt.id
		t.Run(testname, func(t *testing.T) {
			got, got2 := GetJob(tt.id)
			if got2 != tt.want2 {
				t.Errorf("GetJob() = %v, want %v", got2, tt.want2)
			}
			if reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetJob() = %v, want %v", got, tt.want)
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
	testError := UpdateJob(j)

	for _, tt := range tests {
		testname := tt.id
		t.Run(testname, func(t *testing.T) {
			var got models.Job
			if DB.Where("id = ?", tt.id).Scan(&got); reflect.DeepEqual(got, tt.want) {
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
			if got, err := GetJobWithFullHistory(tt.id); reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetJobWithFullHistory() = %v, want %v, error: %v", got, tt.want, err)
			}
		})
	}
}
