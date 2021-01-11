package storage

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/chabad360/covey/models"
	"github.com/chabad360/covey/test"
)

var (
	j  = test.J1
	j2 = test.J2
)

func TestAddJob(t *testing.T) {
	//revive:disable:line-length-limit
	var tests = []struct {
		id   string
		want *models.Job
	}{
		{"update", &j},
		{"3", &models.Job{}},
	}
	//revive:enable:line-length-limit

	testError := AddJob(&j)

	for _, tt := range tests {
		testname := tt.id
		t.Run(testname, func(t *testing.T) {
			var got models.Job
			if err := DB.Where("name = ?", tt.id).First(&got).Error; !cmp.Equal(&got, tt.want) && err != nil {
				t.Errorf("AddJob() = %v, want %v, error: %v", &got, tt.want, testError)
			}
		})
	}
}

func TestGetJob(t *testing.T) {
	DB.Delete(&models.Job{}, "id != ''")
	AddJob(&j)

	//revive:disable:line-length-limit
	var tests = []struct {
		name  string
		id    string
		want  *models.Job
		want2 bool
	}{
		{"success", "update", &j, true},
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
	DB.Delete(&models.Job{}, "id != ''")
	AddJob(&j)
	//revive:disable:line-length-limit
	var tests = []struct {
		name string
		id   string
		want *models.Job
	}{
		{"success", "update", &j},
		{"fail", "3", &models.Job{}},
	}
	//revive:enable:line-length-limit

	j.Cron = "5 * * * *"
	testError := UpdateJob(&j)

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
	DB.Delete(&models.Job{}, "id != ''")
	DB.Delete(&models.Task{}, "id != ''")
	z := task
	ta := &z
	AddTask(ta)
	AddJob(&j)
	ta.IDShort = ""
	ta.Details = nil
	ta.Log = nil
	j.TaskHistory = append(j.TaskHistory, ta.ID)
	UpdateJob(&j)

	jw := &models.JobWithTasks{}
	jw.ID = j.ID
	jw.IDShort = j.GetIDShort()
	jw.Name = j.Name
	jw.Nodes = j.Nodes
	jw.Tasks = j.Tasks
	jw.Cron = j.Cron
	jw.CreatedAt = j.CreatedAt
	jw.UpdatedAt = j.UpdatedAt
	jw.TaskHistory = append(jw.TaskHistory, *ta)

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
