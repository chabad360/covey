package storage

import (
	"github.com/chabad360/covey/models"
	"github.com/google/go-cmp/cmp"
	"testing"
)

var (
	j = &models.Job{
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
	j2 = &models.Job{
		Name:  "add",
		ID:    "3748ffc302b6920c2589795ed6a7cad067eb8f8cb31b079725d0a20bfe6c3b6e",
		Nodes: []string{"node1"},
		Tasks: map[string]models.JobTask{
			"update": {
				Plugin:  "shell",
				Details: map[string]string{"command": "sudo apt update && sudo apt upgrade -y"},
			},
		},
	}
)

func TestAddJob(t *testing.T) {
	//revive:disable:line-length-limit
	var tests = []struct {
		id   string
		want *models.Job
	}{
		{"update", j},
		{"3", &models.Job{}},
	}
	//revive:enable:line-length-limit

	testError := AddJob(j)

	for _, tt := range tests {
		testname := tt.id
		t.Run(testname, func(t *testing.T) {
			var got models.Job
			if err := DB.Where("name = ?", tt.id).First(&got).Error; !cmp.Equal(got, tt.want) && err != nil {
				t.Errorf("AddJob() = %v, want %v, error: %v", got, tt.want, testError)
			}
		})
	}
}

func TestGetJob(t *testing.T) {
	//revive:disable:line-length-limit
	var tests = []struct {
		name  string
		id    string
		want  *models.Job
		want2 bool
	}{
		{"success", "update", j, true},
		{"fail", "3", nil, false},
	}
	//revive:enable:line-length-limit

	for _, tt := range tests {
		testname := tt.name
		t.Run(testname, func(t *testing.T) {
			got, got2 := GetJob(tt.id)
			if got2 != tt.want2 {
				t.Errorf("GetJob() = %v, want %v", got2, tt.want2)
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("GetJob() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateJob(t *testing.T) {
	//revive:disable:line-length-limit
	var tests = []struct {
		name string
		id   string
		want *models.Job
	}{
		{"success", "update", j},
		{"fail", "3", &models.Job{}},
	}
	//revive:enable:line-length-limit

	j.Cron = "5 * * * *"
	testError := UpdateJob(j)

	for _, tt := range tests {
		testname := tt.name
		t.Run(testname, func(t *testing.T) {
			var got models.Job
			if DB.Where("name = ?", tt.id).Scan(&got); !cmp.Equal(&got, tt.want) {
				t.Errorf("UpdateJob() = %v, want %v, error: %v", got, tt.want, testError)
			}
		})
	}
}

func TestGetJobWithFullHistory(t *testing.T) {
	jw := &models.JobWithTasks{}
	jw.Job = *j

	//revive:disable:line-length-limit
	var tests = []struct {
		name  string
		id    string
		want  *models.JobWithTasks
		want2 bool
	}{
		{"success", "update", jw, true},
		{"fail", "3", nil, false},
	}
	//revive:enable:line-length-limit

	for _, tt := range tests {
		testname := tt.name
		t.Run(testname, func(t *testing.T) {
			got, ok := GetJobWithFullHistory(tt.id)

			if ok != tt.want2 {
				t.Errorf("GetJobWithFullHistory() = %v, want %v", ok, tt.want2)
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("GetJobWithFullHistory() = %v, want %v", got, tt.want)
			}
		})
	}
}
